package main

import (
	"io"
	"log"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func initLoggers(loggerOut io.Writer) {
	DebugLogger = log.New(loggerOut, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(loggerOut, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(loggerOut, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func checkError(e error) {
	if e != nil {
		ErrorLogger.Fatalln(e)
	}
}

// Remove trailing slash if necessary
func removeTrailSlash(p string) string {
	if p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}
	return p
}
