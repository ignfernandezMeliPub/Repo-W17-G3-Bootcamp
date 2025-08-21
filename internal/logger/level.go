package logger

type LogLevel int

const (
	LogLevelNone  LogLevel = 0
	LogLevelError LogLevel = 1
	LogLevelAudit LogLevel = 2
	LogLevelDebug LogLevel = 3
)
