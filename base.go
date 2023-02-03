package main

import (
	"io"
	"log"
)

const (
	DMRoot           = ".dm"
	DMConfigFile     = "config.ini"
	DatasetFile      = "datasets.ini"
	GlobalIgnoreFile = "ignore.txt"
)

const GlobalIgnoreDefault = `
# Created by https://www.toptal.com/developers/gitignore/api/macos
# Edit at https://www.toptal.com/developers/gitignore?templates=macos

### macOS ###
# General
.DS_Store
.AppleDouble
.LSOverride

# Icon must end with two \r
Icon


# Thumbnails
._*

# Files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

### macOS Patch ###
# iCloud generated files
*.icloud

# End of https://www.toptal.com/developers/gitignore/api/macos
`

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

const (
	LevelDebug   = 40
	LevelInfo    = 30
	LevelWarning = 20
	LevelError   = 10
	LevelOff     = 0
)

func initLoggers(loggerOut io.Writer, logLevel int) {
	const logFlags int = log.Ldate | log.Ltime | log.Lshortfile

	DebugLogger = log.New(loggerOut, "DEBUG: ", logFlags)
	InfoLogger = log.New(loggerOut, "INFO: ", logFlags)
	WarnLogger = log.New(loggerOut, "WARNING: ", logFlags)
	ErrorLogger = log.New(loggerOut, "ERROR: ", logFlags)

	if logLevel < LevelDebug {
		DebugLogger.SetOutput(io.Discard)
	}
	if logLevel < LevelInfo {
		InfoLogger.SetOutput(io.Discard)
	}
	if logLevel < LevelWarning {
		WarnLogger.SetOutput(io.Discard)
	}
	if logLevel < LevelError {
		ErrorLogger.SetOutput(io.Discard)
	}
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
