package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/posener/complete"
	"github.com/willabides/kongplete"
	"os"
	"os/exec"
	"path/filepath"
)

// CLI Main.
// Main command
var cli struct {
	Info               cliInfo                      `cmd:"" help:"Show infos about this repository."`
	Edit               cliEdit                      `cmd:"" help:"Edit config files."`
	List               cliList                      `cmd:"" help:"List available remote files or directories."`
	Push               cliPullPush                  `cmd:"" help:"Push local state of dataset to remote repository."`
	Pull               cliPullPush                  `cmd:"" help:"Pull remote state of dataset to local."`
	InstallCompletions kongplete.InstallCompletions `cmd:"" help:"install shell completions"`
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

	// Dry run flag is opposite of confirm flag
	ctx.rc.DryRun = !r.Confirm

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

// CLI sub command.
// Show infos.
type cliInfo struct{}

func (r *cliInfo) Run(ctx *cliContext) error {
	// General info
	fmt.Println("Dataset Manager 2")
	fmt.Println("-----------------")
	fmt.Println("Root directory: " + DMRoot)
	fmt.Println("-> Config file:         ", DMConfigFile)
	fmt.Println("-> Dataset definitions: ", DatasetFile)
	fmt.Println("-> Global ignore list:  ", GlobalIgnoreFile)
	fmt.Println()

	// Dataset info
	fmt.Println("Configured datasets:")
	for _, d := range getDatasetNames(ctx.datasets) {
		fmt.Println("* " + d)
	}

	return nil
}

// CLI sub command.
// Edit config files
type cliEdit struct {
	File string `arg:"" default:"datasets" enum:"config,datasets,ignore"`
}

func (r *cliEdit) Run(ctx *cliContext) error {
	// Selct which file to edit
	var file string
	switch r.File {
	case "config":
		file = DMConfigFile
	case "ignore":
		file = GlobalIgnoreFile
	default:
		file = DatasetFile
	}

	// Prepare vim for editing
	cmd := exec.Command("vim", filepath.Join(DMRoot, file))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	// Run it
	err := cmd.Run()
	checkError(err)

	return nil
}

func main() {
	// Init loggers
	initLoggers(os.Stdout, LevelError)

	// Setup CLI parser and bash completer
	cliParser := kong.Must(&cli)
	kongplete.Complete(cliParser,
		kongplete.WithPredictor("ds", complete.PredictDirs("./")),
	)

	// Parse arguments
	ctx, err := cliParser.Parse(os.Args[1:])
	cliParser.FatalIfErrorf(err)

	// Check if root exists
	_, err = os.Stat(DMRoot)
	if os.IsNotExist(err) {
		ErrorLogger.Fatalf("This directory does not contain the `%s` root directory. Please create it manually.", DMRoot)
	}

	// Setup rclone wrapper
	rc := NewRcloneWrapper(filepath.Join(DMRoot, DMConfigFile), true)

	// Load global ignores
	ignores := loadGlobalIgnore(filepath.Join(DMRoot, GlobalIgnoreFile))

	// Load datasets from file and prefix with rclone basedir
	datasets := loadDatasetConfig(filepath.Join(DMRoot, DatasetFile), ignores)

	// Run cli
	err = ctx.Run(&cliContext{
		rc:       rc,
		remote:   rc.preparePath(""),
		datasets: datasets,
		args:     ctx.Args,
	})
	ctx.FatalIfErrorf(err)
}
