// +build !js,!wasm

package QuantoError

import (
	"fmt"
	"github.com/graphql-go/graphql/gqlerrors"
	"log"
	"runtime/debug"
)

var stackEnabled = true

func EnableStackTrace() {
	stackEnabled = true
}

func DisableStackTrace() {
	stackEnabled = false
}

//Flag to define if a stack trace is returned in response or not
func ShowStackTrace() bool {
	return stackEnabled
}

type ErrorObject struct {
	ErrorCode  string      `json:"errorCode"`
	ErrorField string      `json:"errorField"`
	Message    string      `json:"message"`
	ErrorData  interface{} `json:"errorData"`
	StackTrace string      `json:"stackTrace"`
}

func New(errorCode, errorField, message string, errorData interface{}) *ErrorObject {
	return &ErrorObject{
		ErrorCode:  errorCode,
		ErrorField: errorField,
		ErrorData:  errorData,
		Message:    message,
		StackTrace: string(debug.Stack()),
	}
}

func (e *ErrorObject) Error() string {
	return e.Message
}

func (e *ErrorObject) ToFormattedError() gqlerrors.FormattedError {
	log.Println(e)
	baseErr := gqlerrors.FormatError(e)

	if baseErr.Extensions == nil {
		baseErr.Extensions = make(map[string]interface{})
	}

	//baseErr.Extensions["errorObject"] = e
	baseErr.Extensions["errorCode"] = e.ErrorCode
	baseErr.Extensions["errorField"] = e.ErrorField
	if stackEnabled {
		baseErr.Extensions["stackTrace"] = e.StackTrace
	}
	baseErr.Extensions["errorData"] = e.ErrorData

	return baseErr
}

func (e *ErrorObject) String() string {
	o := fmt.Sprintf("Error: %s\n", e.Message)
	o += fmt.Sprintf("  Error Code: %s\n", e.ErrorCode)
	o += fmt.Sprintf("  Error Field: %s\n", e.ErrorField)
	o += fmt.Sprintf("  Error Data: %v\n", e.ErrorData)
	o += fmt.Sprintf("  Stack Trace %s\n", e.StackTrace)
	return o
}
