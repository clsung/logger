package logger

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
)

type LogFunc func(message string)

const OUTFILE = "out.json"

func createOutFile() *os.File {
	// Delete file first if exists
	os.Remove(OUTFILE)

	file, err := os.Create(OUTFILE)
	if err != nil {
		panic("Unable to create test file")
	}

	return file
}

func compareWithOutFile(expected string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	return strings.TrimRight(string(data), "\n") == expected
}

func outFileContains(substring string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	fileData := strings.TrimRight(string(data), "\n")
	return strings.Contains(fileData, substring)
}

func TestLoggerDebug(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New("robokiller-ivr", "1.0")
	log.SetWriter(file)

	log.Set("key", "value")
	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"data\":{\"key\":\"value\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Fail()
	}
}

func TestLoggerInfo(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New("robokiller-ivr", "1.0")
	log.SetWriter(file)

	log.Set("key", "value")
	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"data\":{\"key\":\"value\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Fail()
	}
}

func TestLoggerError(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	log := New("robokiller-ivr", "1.0")
	log.SetWriter(file)

	log.Set("key", "value")
	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"reportLocation\":{\"filePath\":", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not a stacktrace key", OUTFILE)
	}
}
