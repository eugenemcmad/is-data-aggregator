package utils

import (
	"fmt"
)

const PanicErrorString = "%v PANIC:'%v'\n\tat:[[\n%s\n]]\n"

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Create panic error.
func PanicErr(irecover interface{}, stack []byte, method string, params ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = Errorf("err-helper.PanicErr(recover:%v, stack:%v, method:%v, params:[%v]) PROCESSING CRASH!", irecover, stack, method, params)
		}
	}()

	return Errorf(PanicErrStr(irecover, stack, method, params)) // _note: go magic test(a,b)-->test([a,b])
}

// Create panic error description.
func PanicErrStr(irecover interface{}, stack []byte, method string, params ...interface{}) (errStr string) {
	defer func() {
		if r := recover(); r != nil {
			errStr = fmt.Sprintf("err-helper.PanicErrStr(recover:%v, stack:%v, method:%v, params:[%v]) PROCESSING CRASH!", irecover, stack, method, params)
		}
	}()

	method = getMethodDesc(method, params...)

	errStr = fmt.Sprintf(PanicErrorString, method, irecover, stack)
	return errStr
}

func getMethodDesc(method string, params ...interface{}) (methodDesc string) {

	defer func() {
		if r := recover(); r != nil {
			methodDesc = fmt.Sprintf("err-helper.getMethodDesc(method:%v, params:[%v]) PROCESSING CRASH!", method, params)
		}
	}()

	len := len(params)
	if method == "" {
		method = "uncknown_method_"
	}
	method = method + "("
	for i := 0; i < len; i++ {
		method += fmt.Sprintf("%v", params[i])
		if i < (len - 1) {
			method += ", "
		}
	}
	method = method + ")"

	return method
}
