package main

import (
	"github.com/alecthomas/kong"
	"os"
)

// CLI Main.
// Main command
var cli struct {
	List cliList     `cmd:"" help:"List available remote files or directories."`
	Push cliPullPush `cmd:"" help:"Push local state of dataset to remote repository."`
	Pull cliPullPush `cmd:"" help:"Pull remote state of dataset to local."`
}

// CLI context.
// Context provided to Run(...) method of CLI commands
type cliContext struct {
	rc       *RcloneWrapper
	remote   string
	datasets map[string]*Dataset
	args     []string
}

// Return a Dataset pointer for given dataset name
func (ctx *cliContext) getDataset(dsName string) (*Dataset, string) {
	// Get dataset from map
	ds := ctx.datasets[dsName]
	if ds == nil {
		ErrorLogger.Fatalf("Unknown dataset %s", dsName)
	}

	// Create temporary filter file because all CLI commands need it
	filterPath := ds.WriteTmpFilterFile()

	return ds, filterPath
}

// CLI sub command.
// Push and pull datasets.
type cliPullPush struct {
	Ds      string `arg:"" help:"Dataset name."`
	Confirm bool   `help:"Add flag to confirm action."`
}

func (r *cliPullPush) Run(ctx *cliContext) error {
	_, filterPath := ctx.getDataset(r.Ds)

	switch ctx.args[0] {
	case "pull":
		ctx.rc.exec("sync", ctx.remote, "--filter-from", filterPath, ".")
	case "push":
		ctx.rc.exec("sync", ".", "--filter-from", filterPath, ctx.remote)
	}

	err := os.Remove(filterPath)
	checkError(err)

	return nil
}

// CLI sub command.
// List remote directories
type cliList struct {
	Ds string `arg:"" optional:"" help:"Dataset to list (optional). If not provided, all available directories are listed."`
}

func (r *cliList) Run(ctx *cliContext) error {
	if r.Ds == "" {
		// List directories
		ctx.rc.List()
	} else {
		// List files of requested dataset
		_, filterPath := ctx.getDataset(r.Ds)
		ctx.rc.exec("ls", ctx.remote, "--filter-from", filterPath)

		err := os.Remove(filterPath)
		checkError(err)
	}
	return nil
}

func main() {
	// Init loggers
	initLoggers(os.Stdout, LevelError)

	// Setup rclone wrapper
	rc := NewRcloneWrapper("example/dm.config.ini", true)

	// Load global ignores
	ignores := loadGlobalIgnore("example/global_ignore.txt")

	// Load datasets from file and prefix with rclone basedir
	datasets := loadDatasetConfig("example/datasets.example.ini", ignores)

	// Parse arguments and run CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&cliContext{
		rc:       rc,
		remote:   rc.preparePath(""),
		datasets: datasets,
		args:     ctx.Args,
	})
	ctx.FatalIfErrorf(err)
}
