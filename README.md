# teltech/log

Super simple structured logging mechanism for Go projects with [Stackdriver format](https://cloud.google.com/error-reporting/docs/formatting-error-messages) compatibility

## Installation

``` sh
go get -u github.com/teltech/log
```

## Usage
``` go
package main

import (
    "github.com/teltech/log"
)

// There should be a LOG_LEVEL environment variable set, which is read by the library
// If no value is set, the default LOG_LEVEL will be INFO

func main() {
    logger, err := log.New("project-name", "project-version")
    if err != nil {
        fmt.Print("Cannot initiate logger")
    }

    // A metric is an INFO log entry without a payload
    logger.Metric("CUSTOM_METRIC_ENTRY")

    logger.Set("user", "+1234567890")
    logger.Set("action", "create-account")

    // Log a DEBUG message, only visible in when LOG_LEVEL is set to DEBUG
    logger.Debug("debug message goes here")

    // Log an INFO message
    logger.Info("info message goes here")

    // Log a WARN message
    logger.Warn("warn message goes here")

    // Error() prints the stacktrace as part of the payload for each entry and sends the
    // data to Stackdriver Error Reporting service
    logger.Error("error message goes here")
}