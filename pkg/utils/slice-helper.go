package utils

import "fmt"

func GetMaxValue(Data []int) (int, error) {
	if len(Data) == 0 {
		return 0, fmt.Errorf("slice is empty")
	}

	maxVal := Data[0]

	for _, value := range Data {
		if value > maxVal {
			maxVal = value
		}
	}

	return maxVal, nil
}
