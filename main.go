package logger

import (
	"time"
	"fmt"
	"encoding/json"
	"runtime"
	"bytes"
	"os"
	"io"
)


const (
	LOG_LEVEL_DEBUG = 0
	LOG_LEVEL_INFO  = 1
	LOG_LEVEL_WARN  = 2
	LOG_LEVEL_ERROR = 3
)

var LogLevelName = map[int]string{
	0: "DEBUG",
	1: "INFO",
	2: "WARN",
	3: "ERROR",
}

var LogLevelValue = map[string]int{
	"DEBUG": 0,
	"INFO": 1,
	"WARN": 2,
	"ERROR": 3,
}

type ServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

type ReportLocation struct {
	FilePath     string `json:"filePath"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
}

type Context struct {
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

type Payload struct {
	Severity       string             `json:"severity"`
	EventTime      string             `json:"eventTime"`
	Caller         string             `json:"caller,omitempty"`
	Message        string             `json:"message"`
	Data           map[string]string  `json:"data,omitempty"`
	ServiceContext *ServiceContext    `json:"serviceContext"`
	Context        *Context           `json:"context,omitempty"`
	Stacktrace     string             `json:"stacktrace,omitempty"`
}

type Log struct {
	Payload *Payload
	writer io.Writer
}

func New(service, version string) *Log {
	return &Log{
		Payload: &Payload{
			ServiceContext: &ServiceContext{
				Service: service,
				Version: version,
			},
		},
		writer: os.Stdout,
	}
}

func (l *Log) SetWriter(w io.Writer) {
	l.writer = w
}

func (l *Log) Set(key, val string) {
	if l.Payload.Data == nil {
		l.Payload.Data = map[string]string{}
	}

	l.Payload.Data[key] = val
}

func (l *Log) log(severity, message string) {
	l.Payload = &Payload{
		Severity: severity,
		EventTime: time.Now().Format(time.RFC3339),
		Message: message,
		Data: l.Payload.Data,
		ServiceContext: l.Payload.ServiceContext,
		Context: l.Payload.Context,
		Stacktrace: l.Payload.Stacktrace,
	}

	line, ok := json.Marshal(l.Payload)
	if ok != nil {
		fmt.Errorf("cannot marshal payload: %s", ok.Error())
	}

	fmt.Fprintln(l.writer, string(line))

	// Unset the current payload data
	l.Payload.Data = nil
}

// Checks whether the specified log level is valid in the current environment
func isValidLogLevel(logLevel int) bool {
	curLogLev, ok := LogLevelValue[os.Getenv("LOG_LEVEL")]
	if !ok {
		fmt.Errorf("the LOG_LEVEL environment variable is not set or has an incorrect value")
	}

	return curLogLev <= logLevel
}

func (l *Log) Debug(message string) {
	if !isValidLogLevel(LOG_LEVEL_DEBUG) {
		return
	}

	l.log(LogLevelName[LOG_LEVEL_DEBUG], message)
}

func (l *Log) Info(message string) {
	if !isValidLogLevel(LOG_LEVEL_INFO) {
		return
	}

	l.log(LogLevelName[LOG_LEVEL_INFO], message)
}

func (l *Log) Warn(message string) {
	if !isValidLogLevel(LOG_LEVEL_WARN) {
		return
	}

	l.log(LogLevelName[LOG_LEVEL_WARN], message)
}

func (l *Log) Error(message string) {
	buffer := make([]byte, 1024)
	runtime.Stack(buffer, false)
	_, file, line, _ := runtime.Caller(1)
	l.Payload = &Payload{
		ServiceContext: l.Payload.ServiceContext,
		Context: &Context{
			ReportLocation: &ReportLocation{
				FilePath: file,
				FunctionName: "unknown",
				LineNumber: line,
			},
		},
		Stacktrace: string(bytes.Trim(buffer, "\x00")),
	}

	l.log(LogLevelName[LOG_LEVEL_ERROR], message)
}