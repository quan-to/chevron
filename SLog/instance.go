package SLog

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
)

const logBaseFormat = "%40v | %s"

type Instance struct {
	scope string
}

func (i *Instance) LogNoFormat(str interface{}, v ...interface{}) *Instance {
	if infoEnabled {
		var c = aurora.Cyan
		log.Printf(logBaseFormat, c(aurora.Bold(i.scope)).String(), fmt.Sprintf(asString(str), v...))
	}
	return i
}

func (i *Instance) Log(str interface{}, v ...interface{}) *Instance {
	if infoEnabled {
		var c = aurora.Cyan
		log.Printf(logBaseFormat, c(aurora.Bold(i.scope)).String(), c(fmt.Sprintf(asString(str), v...)))
	}
	return i
}

func (i *Instance) Info(str interface{}, v ...interface{}) *Instance {
	return i.Log(str, v...)
}

func (i *Instance) Debug(str interface{}, v ...interface{}) *Instance {
	if debugEnabled {
		var c = aurora.Magenta
		log.Printf(logBaseFormat, c(aurora.Bold(i.scope)).String(), c(fmt.Sprintf(asString(str), v...)))
	}
	return i
}

func (i *Instance) Warn(str interface{}, v ...interface{}) *Instance {
	if warnEnabled {
		var c = aurora.Brown
		log.Printf(logBaseFormat, c(aurora.Bold(i.scope)).String(), c(fmt.Sprintf(asString(str), v...)))
	}
	return i
}

func (i *Instance) Error(str interface{}, v ...interface{}) *Instance {
	if errorEnabled {
		var c = aurora.Red
		log.Printf(logBaseFormat, c(aurora.Bold(i.scope)).String(), c(fmt.Sprintf(asString(str), v...)))
	}
	return i
}

func (i *Instance) Fatal(str interface{}, v ...interface{}) {
	var msg = fmt.Sprintf(asString(str), v...)
	i.Error(msg)
	panic(msg)
}
