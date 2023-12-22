package log


// Print to terminal
func NewLoggerPrint() LogCB {
	return func() (logger DoLog, closer Closer) {
		logger = func (lvl LogLevel, prefix, msg string) {
			GetColor(lvl).Println(prefix + " " + msg)
		}
		return
	}
}
