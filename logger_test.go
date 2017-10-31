package logger

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
)

const outfile = "out.json"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	initConfig(debug, "robokiller-ivr", "1.0")
}

func createOutFile() *os.File {
	// Delete file first if exists
	os.Remove(outfile)

	file, err := os.Create(outfile)
	if err != nil {
		panic("Unable to create test file")
	}

	return file
}

func compareWithOutFile(expected string) bool {
	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		panic("Unable to read test file")
	}

	return strings.TrimRight(string(data), "\n") == expected
}

func outFileContains(substring string) bool {
	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		panic("Unable to read test file")
	}

	fileData := strings.TrimRight(string(data), "\n")
	return strings.Contains(fileData, substring)
}

func outFileDoesNotContain(substring string) bool {
	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		panic("Unable to read test file")
	}

	fileData := strings.TrimRight(string(data), "\n")
	return !strings.Contains(fileData, substring)
}

func TestLoggerDebugWithImplicitContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"key": "value",
		"function" : "TestLoggerDebug",
	}).SetWriter(file)

	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerDebugWithoutContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().SetWriter(file)

	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerDebugfWithoutContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().SetWriter(file)

	param := "with param"
	log.Debugf("debug message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerMetric(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().SetWriter(file)

	log.Metric("custom_metric")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"custom_metric\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerInfo(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerInfo",
	}).SetWriter(file)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerInfof(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerInfo",
	}).SetWriter(file)

	param := "with param"
	log.Infof("info message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerError(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerError",
	}).SetWriter(file)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", outfile, expected)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"}") {
		t.Errorf("output file %s does not contain the context", outfile)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", outfile)
	}
}

func TestLoggerErrorWithoutContext(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().SetWriter(file)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"reportLocation\"", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", outfile, expected)
	}

	// Check that the error entry contains the context
	if !outFileDoesNotContain("\"context\":{\"data\":") {
		t.Errorf("output file %s has a context nad it wasn't supposed to", outfile)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", outfile)
	}
}

func TestLoggerErrorf(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"key":"value",
		"function": "TestLoggerError",
	}).SetWriter(file)

	param := "with param"
	log.Errorf("error message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", outfile, expected)
	}
}

func TestLoggerInfoWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"function": "TestLoggerInfo",
		"key":"value",
		"package": "logger",
	}).SetWriter(file)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\",\"package\":\"logger\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", outfile, expected)
	}
}

func TestLoggerErrorWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New().With(Fields{
		"function": "TestLoggerError",
		"key":"value",
		"package": "logger",
	}).SetWriter(file)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", outfile, expected)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\",\"package\":\"logger\"}") {
		t.Errorf("output file %s does not contain the context", outfile)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", outfile)
	}
}
