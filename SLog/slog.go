package SLog

// TODO: Syslog Output

// region Global
var debugEnabled = true
var warnEnabled = true
var errorEnabled = true
var infoEnabled = true

var glog = Instance{scope: "RemoteSigner"}

func LogNoFormat(str interface{}, v ...interface{}) *Instance {
	return glog.LogNoFormat(str, v...)
}

func Log(str interface{}, v ...interface{}) *Instance {
	return glog.Log(str, v...)
}

func Info(str interface{}, v ...interface{}) *Instance {
	return glog.Info(str, v...)
}

func Debug(str interface{}, v ...interface{}) *Instance {
	return glog.Debug(str, v...)
}

func Warn(str interface{}, v ...interface{}) *Instance {
	return glog.Warn(str, v...)
}

func Error(str interface{}, v ...interface{}) *Instance {
	return glog.Error(str, v...)
}

func Fatal(str interface{}, v ...interface{}) {
	glog.Fatal(str, v)
}

func Scope(scope string) *Instance {
	return &Instance{
		scope: scope,
	}
}

func SetDebug(enabled bool) {
	debugEnabled = enabled
}
func SetWarning(enabled bool) {
	warnEnabled = enabled
}
func SetInfo(enabled bool) {
	infoEnabled = enabled
}
func SetError(enabled bool) {
	errorEnabled = enabled
}

func SetTestMode() {
	SetDebug(false)
	SetWarning(false)
	SetInfo(false)
	SetError(false)
}

func UnsetTestMode() {
	SetDebug(true)
	SetWarning(true)
	SetInfo(true)
	SetError(true)
}

func DebugEnabled() bool {
	return debugEnabled
}
func WarningEnabled() bool {
	return warnEnabled
}
func InfoEnabled() bool {
	return infoEnabled
}
func ErrorEnabled() bool {
	return errorEnabled
}

// endregion
