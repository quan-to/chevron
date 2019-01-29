package QuantoError

import (
	"github.com/quan-to/graphql/gqlerrors"
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
