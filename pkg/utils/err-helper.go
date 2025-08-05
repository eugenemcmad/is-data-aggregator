// Package utils provides utility functions for error handling and data manipulation.
// This package contains helper functions for creating formatted errors, handling panics,
// and managing slice operations.
package utils

import (
	"fmt"
)

// PanicErrorString is the format string used for creating panic error messages.
// It includes placeholders for method name, panic value, and stack trace.
const PanicErrorString = "%v PANIC:'%v'\n\tat:[[\n%s\n]]\n"

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// This is a convenience wrapper around fmt.Errorf for creating formatted errors.
//
// Example:
//
//	err := Errorf("failed to process %s", "data")
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// PanicErr creates a panic error with detailed information including the panic value,
// stack trace, method name, and parameters. This function is designed to be used
// in defer statements to capture panic information.
//
// Parameters:
//   - irecover: the panic value (typically from recover())
//   - stack: the stack trace as bytes
//   - method: the name of the method where the panic occurred
//   - params: additional parameters to include in the error message
//
// Returns:
//   - error: a formatted error containing panic details
//
// Example:
//
//	defer func() {
//	    if r := recover(); r != nil {
//	        err = PanicErr(r, debug.Stack(), "ProcessData", data)
//	    }
//	}()
func PanicErr(irecover interface{}, stack []byte, method string, params ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = Errorf("err-helper.PanicErr(recover:%v, stack:%v, method:%v, params:[%v]) PROCESSING CRASH!", irecover, stack, method, params)
		}
	}()

	return Errorf(PanicErrStr(irecover, stack, method, params)) // _note: go magic test(a,b)-->test([a,b])
}

// PanicErrStr creates a panic error description string with detailed information.
// This function formats the panic information into a readable string without
// returning an error type.
//
// Parameters:
//   - irecover: the panic value (typically from recover())
//   - stack: the stack trace as bytes
//   - method: the name of the method where the panic occurred
//   - params: additional parameters to include in the error message
//
// Returns:
//   - string: a formatted string containing panic details
//
// Example:
//
//	panicStr := PanicErrStr(r, debug.Stack(), "ProcessData", data)
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

// getMethodDesc formats a method name with its parameters into a readable string.
// This is an internal helper function used by PanicErr and PanicErrStr.
//
// Parameters:
//   - method: the name of the method
//   - params: the parameters to format with the method name
//
// Returns:
//   - string: a formatted method description like "methodName(param1, param2)"
//
// Example:
//
//	desc := getMethodDesc("ProcessData", "input", 42)
//	// Returns: "ProcessData(input, 42)"
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
