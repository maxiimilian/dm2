package main

import (
	"fmt"
	"os"
)

func listFiles(rc *RcloneWrapper, ds *Dataset) {
	filterPath := ds.WriteTmpFilterFile()
	remote := rc.preparePath("")
	rc.exec("ls", remote, "--filter-from", filterPath)
}

func pullDataset(rc *RcloneWrapper, ds *Dataset) {

}

func pushDataset(rc *RcloneWrapper, ds *Dataset) {

}

func main() {
	// Init loggers
	initLoggers(os.Stdout)

	// Setup rclone wrapper
	rc := NewRcloneWrapper("example/dm.config.ini", true)
	fmt.Println(rc)

	// Load global ignores
	ignores := loadGlobalIgnore("example/global_ignore.txt")

	// Load datasets from file and prefix with rclone basedir
	datasets := loadDatasetConfig("example/datasets.example.ini", ignores)

	listFiles(rc, datasets[2])
}
