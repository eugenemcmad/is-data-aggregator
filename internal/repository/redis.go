package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
	"xis-data-aggregator/internal/api"
	"xis-data-aggregator/pb"

	"github.com/redis/go-redis/v9"

	"xis-data-aggregator/internal/models"
)

const (
	zsetKey = "events"
	ttlSec  = 500
)

var (
	ErrNotFound = errors.New("not found")
	ErrCorrupt  = errors.New("corrupted data")
)

var ctx = context.Background()

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository() (*RedisRepository, error) {
	repo := RedisRepository{}
	err := repo.Open()
	return &repo, err
}

func (o *RedisRepository) Open() error {
	srv, err := miniredis.Run()
	if err != nil {
		return err
	}
	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	o.Client = client

	return err
}

// Close closes database connections.
func (o *RedisRepository) Close() error {
	return o.Client.Close()
}

func (o *RedisRepository) Put(data *models.Data) error {

	pbData, err := api.DataToProto(data)
	switch {
	case err != nil:
		return err
	case pbData == nil:
		return fmt.Errorf("pbData is nil")
	}

	bytes, err := proto.Marshal(pbData)
	if err != nil {
		return err
	}

	// Time range `table` without TTL. Partitioning is recommended, by month for example
	_, err = o.Client.ZAdd(ctx, zsetKey, redis.Z{Score: float64(pbData.Timestamp), Member: bytes}).Result()
	if err != nil {
		return err
	}

	// Fast key-value `table` with TTL
	err = o.Client.Set(ctx, pbData.Id, bytes, ttlSec*time.Second).Err()

	return err
}

func (o *RedisRepository) GetByID(id uuid.UUID) (*models.Data, error) {
	val, err := o.Client.Get(ctx, id.String()).Bytes()

	switch {
	case errors.Is(err, redis.Nil):
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	}

	var umData pb.Data
	err = proto.Unmarshal(val, &umData)
	if err != nil {
		return nil, err // todo:  corrupted data
	}

	return api.ProtoToData(&umData)

}

func (o *RedisRepository) ListByPeriod(from, to int64) ([]models.Data, error) {
	var res []models.Data

	results, err := o.Client.ZRangeByScoreWithScores(ctx, zsetKey, &redis.ZRangeBy{
		Min: strconv.Itoa(int(from)),
		Max: strconv.Itoa(int(to)),
	}).Result()

	switch {
	case err != nil:
		return nil, err
	case results == nil:
		return nil, ErrNotFound
	case len(results) == 0:
		return nil, ErrNotFound
	}

	// Note: any invalid entry compromises the entire set -> 1 error -> return
	for _, result := range results {
		var umData pb.Data
		err = proto.Unmarshal(result.Member.([]byte), &umData)
		if err != nil {
			return nil, ErrCorrupt
		}
		data, err := api.ProtoToData(&umData)
		if err != nil {
			return nil, err
		}
		res = append(res, *data)
	}

	return res, nil
}

/*Общий Принцип и Рекомендации
Repository: Должен быть источником истины о том, что объект не найден. Он должен возвращать nil для объекта и специальную, экспортируемую ошибку (например, repository.ErrNotFound). Используйте errors.Is для проверки этой ошибки.

Service: Принимает решение о том, как интерпретировать ошибку "не найдено" от репозитория.

Если "не найдено" является ожидаемым результатом (например, GET /users/{id} и пользователя нет), то Service должен просто передать (или обернуть) ошибку repository.ErrNotFound выше, возможно, преобразуя ее в ошибку уровня Service (service.ErrUserNotFound), чтобы сохранить абстракцию.
Если "не найдено" является бизнес-исключением (например, при попытке обновить несуществующий ресурс, или если связанный ресурс не найден), Service должен сгенерировать соответствующую ошибку бизнес-логики (например, ErrInvalidInput, ErrPreconditionFailed).

HTTP-обработчик (или RPC-слой): Отвечает за окончательное преобразование ошибок из сервисного слоя в соответствующий HTTP-статус (например, repository.ErrNotFound или service.ErrUserNotFound -> 404 Not Found).
В большинстве случаев, когда вы просто хотите получить сущность по ID и сообщить, что ее нет, я бы рекомендовал:

Repository возвращает nil, repository.ErrNotFound.

Service просто возвращает nil, service.ErrUserNotFound (используя errors.Is для проверки repository.ErrNotFound и оборачивая/преобразуя ее).

HTTP-обработчик проверяет errors.Is(err, service.ErrUserNotFound) и возвращает 404.
*/
