package logger

type LogLayer string

var (
	LogLayerHandler    LogLayer = "handler"
	LogLayerService    LogLayer = "service"
	LogLayerRepository LogLayer = "repository"
)
