# teltech/logger

Super simple structured logging mechanism for Go projects with [Stackdriver format](https://cloud.google.com/error-reporting/docs/formatting-error-messages) compatibility

## Installation

``` sh
go get -u github.com/teltech/logger
```

## Usage
``` go
package main

import (
    "github.com/teltech/logger"
)

// There should be a LOG_LEVEL environment variable set, which is read by the library
// If no value is set, the default LOG_LEVEL will be INFO

func main() {
    log, err := log.New("project-name", "project-version")
    if err != nil {
        fmt.Print("Cannot initiate logger")
    }

    // A metric is an INFO log entry without a payload
    log.Metric("CUSTOM_METRIC_ENTRY")

    log.Set("user", "+1234567890")
    log.Set("action", "create-account")

    // Log a DEBUG message, only visible in when LOG_LEVEL is set to DEBUG
    log.Debug("debug message goes here")

    // Log an INFO message
    log.Info("info message goes here")

    // Log a WARN message
    log.Warn("warn message goes here")

    // Error() prints the stacktrace as part of the payload for each entry and sends the
    // data to Stackdriver Error Reporting service
    log.Error("error message goes here")
}