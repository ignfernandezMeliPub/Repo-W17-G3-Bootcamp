package logger

type LogStatus string

const (
	LogStatusInProgress LogStatus = "In Progress"
	LogStatusSuccess    LogStatus = "Success"
	LogStatusError      LogStatus = "Error"
)
