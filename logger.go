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

type severity int

const (
	debug severity = iota
	info
	warn
	error
)

func (s severity)String() string {
	return logLevelName[s]
}

var logLevelName = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

var logLevelValue = map[string]severity{
	"DEBUG": debug,
	"INFO":  info,
	"WARN":  warn,
	"ERROR": error,
}

// Fields is used to wrap the log entries payload
type Fields map[string]string

// ServiceContext is required by the Stackdriver Error format
type ServiceContext struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
}

// ReportLocation is required by the Stackdriver Error format
type ReportLocation struct {
	FilePath     string `json:"filePath"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
}

// Context is required by the Stackdriver Error format
type Context struct {
	Data           Fields          `json:"data,omitempty"`
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

// Payload groups all the data for a log entry
type Payload struct {
	Severity       string          `json:"severity"`
	EventTime      string          `json:"eventTime"`
	Caller         string          `json:"caller,omitempty"`
	Message        string          `json:"message"`
	ServiceContext *ServiceContext `json:"serviceContext"`
	Context        *Context        `json:"context,omitempty"`
	Stacktrace     string          `json:"stacktrace,omitempty"`
}

// Log is the main type for the logger package. The errors require a specific JSON format for them to be
// ingested and processed by Google Cloud Platform Stackdriver Logging and Error Reporting.
// See: https://cloud.google.com/error-reporting/docs/formatting-error-messages
// The resulting output has the following format, optional fields are... well, optional:
// @TODO Move this to the README
/**
    {
       "severity": "ERROR",
       "eventTime": "2017-04-26T02:29:33-04:00",
       "message": "An error just happened!",
       "serviceContext": {
          "service": "robokiller-ivr",
          "version": "1.0"
       },
       "context": {
         "data": {
           "clientIP": "127.0.0.1"
           "userAgent": "Mosaic 1.0"
         },
         "reportLocation": {
           "filePath": "\/Users\/mc\/Documents\/src\/github.com\/macuenca\/apex\/mauricio.go",
           "functionName": "unknown",
           "lineNumber": 15
         }
       },
      "stacktrace": "goroutine 1 [running]:main.main()\n\t\/github.com\/macuenca\/mauricio.go:15 +0x1a9\n"
    }
 */
type Log struct {
	payload *Payload
	writer  io.Writer
}

var (
	logLevel severity
	service  string
	version  string
)

// @TODO Make this call a initConfig() function
func init() {
	ll, ok := logLevelValue[os.Getenv("LOG_LEVEL")]
	if !ok {
		fmt.Println("logger warn: LOG_LEVEL is not valid or not set, defaulting to INFO")
		logLevel = logLevelValue[info.String()]
	} else {
		logLevel = ll
	}

	if os.Getenv("SERVICE") == "" || os.Getenv("VERSION") == "" {
		fmt.Println("logger error: cannot instantiate the logger, make sure the SERVICE and VERSION environment vars are set correctly")
	}

	initConfig(logLevel, os.Getenv("SERVICE"), os.Getenv("VERSION"))
}

func initConfig(lvl severity, svc, ver string) {
	logLevel = lvl
	service = svc
	version = ver
}

// New instantiates and returns a Log object
func New() *Log {
	return &Log{
		payload: &Payload{
			ServiceContext: &ServiceContext{
				Service: service,
				Version: version,
			},
		},
		writer: os.Stdout,
	}
}

// SetWriter exists mainly for tests, allowing to change the output from STDOUT to FILE
func (l *Log) SetWriter(w io.Writer) *Log {
	l.writer = w
	return l
}

func (l *Log) set(key, val string) {
	if l.payload.Context == nil {
		l.payload.Context = &Context{
			Data: Fields{},
		}
	}

	l.payload.Context.Data[key] = val
}

func (l *Log) log(severity, message string) {
	// Do not persist the payload here, just format it, marshal it and return it
	l.payload = &Payload{
		Severity:       severity,
		EventTime:      time.Now().Format(time.RFC3339),
		Message:        message,
		ServiceContext: l.payload.ServiceContext,
		Context:        l.payload.Context,
		Stacktrace:     l.payload.Stacktrace,
	}

	payload, ok := json.Marshal(l.payload)
	if ok != nil {
		fmt.Printf("logger error: cannot marshal payload: %s", ok.Error())
	}

	fmt.Fprintln(l.writer, string(payload))
}

// Checks whether the specified log level is valid in the current environment
func isValidLogLevel(s severity) bool {
	return s >= logLevel
	//curLogLev, ok := logLevelValue[os.Getenv("LOG_LEVEL")]
	//if !ok {
	//	fmt.Println("logger warn: LOG_LEVEL is not set, defaulting to INFO")
	//	os.Setenv("LOG_LEVEL", "INFO")
	//}
	//
	//return curLogLev <= logLevel
}

// With is used as a chained method to specify which values go in the log entry's context
func (l *Log) With(fields Fields) *Log {
	for k, v := range fields {
		l.set(k, v)
	}

	return l
}

// Debug prints out a message with DEBUG severity level
func (l *Log) Debug(message string) {
	if !isValidLogLevel(debug) {
		return
	}

	// @TODO use debug.String()
	l.log(logLevelName[debug], message)
}

// Debugf prints out a message with DEBUG severity level
func (l *Log) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

// Metric prints out a message with INFO severity and no extra fields
func (l *Log) Metric(message string) {
	if !isValidLogLevel(info) {
		return
	}

	l.log(logLevelName[info], message)
}

// Info prints out a message with INFO severity level
func (l *Log) Info(message string) {
	if !isValidLogLevel(info) {
		return
	}

	l.log(logLevelName[info], message)
}

// Infof prints out a message with INFO severity level
func (l *Log) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

// Warn prints out a message with WARN severity level
func (l *Log) Warn(message string) {
	if !isValidLogLevel(warn) {
		return
	}

	l.log(logLevelName[warn], message)
}

// Warnf prints out a message with WARN severity level
func (l *Log) Warnf(message string, args ...interface{}) {
	l.Warn(fmt.Sprintf(message, args...))
}

// Error prints out a message with ERROR severity level
func (l *Log) Error(message string) {
	buffer := make([]byte, 1024)
	runtime.Stack(buffer, false)
	_, file, line, _ := runtime.Caller(1)

	// Set the data when the context is empty
	if l.payload.Context == nil {
		l.payload.Context = &Context{
			Data: Fields{},
		}
	}

	// @TODO Create a new logger here and print it
	l.payload = &Payload{
		ServiceContext: l.payload.ServiceContext,
		Context: &Context{
			Data: l.payload.Context.Data,
			ReportLocation: &ReportLocation{
				FilePath: file,
				FunctionName: "unknown",
				LineNumber: line,
			},
		},
		Stacktrace: string(bytes.Trim(buffer, "\x00")),
	}

	l.log(logLevelName[error], message)
}

// Errorf prints out a message with ERROR severity level
func (l *Log) Errorf(message string, args ...interface{}) {
	l.Error(fmt.Sprintf(message, args...))
}
