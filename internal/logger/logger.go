package logger

import (
	"database/sql"
	"time"
)

var logLevel LogLevel = LogLevelNone
var logDb *sql.DB = nil

type LogInfo struct {
	Source     string
	Endpoint   string
	HttpMethod string
	Layer      LogLayer
	Action     string
	Status     LogStatus
	Message    string
	level      LogLevel
	time       time.Time
}

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetLogDb(db *sql.DB) {
	logDb = db
}

func Debug(logInfo LogInfo) {
	if logLevel < LogLevelDebug {
		return
	}

	logInfo.level = LogLevelDebug
	logInfo.time = time.Now()

	log(logInfo)
}

func Audit(logInfo LogInfo) {
	if logLevel < LogLevelAudit {
		return
	}

	logInfo.level = LogLevelAudit
	logInfo.time = time.Now()

	log(logInfo)
}

func Error(logInfo LogInfo) {
	if logLevel < LogLevelError {
		return
	}

	logInfo.level = LogLevelError
	logInfo.time = time.Now()

	log(logInfo)
}

func log(logInfo LogInfo) {
	if logDb == nil {
		println("logDb is nil")
		return
	}

	_, err := logDb.Exec("INSERT INTO logs (source, endpoint, http_method, layer, action, status, message, level, time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", logInfo.Source, logInfo.Endpoint, logInfo.HttpMethod, logInfo.Layer, logInfo.Action, logInfo.Status, logInfo.Message, logInfo.level, logInfo.time)

	if err != nil {
		println("error inserting log: " + err.Error())
		SetLogDb(nil)
	}
}
