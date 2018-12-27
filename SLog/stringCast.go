package SLog

import (
	"github.com/logrusorgru/aurora"
	"runtime/debug"
)

type StringCast interface {
	String() string
}

func asString(str interface{}) string {
	switch v := str.(type) {
	default:
		debug.PrintStack()
		Fatal(aurora.Red("Unexpected type %T"), v)
		return "" // Linter bug fix
	case StringCast:
		return v.String()
	case error:
		return v.Error()
	case string:
		return v
	}
}
