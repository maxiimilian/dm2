package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"strconv"
)

type RcloneWrapper struct {
	// Static variables read from config file
	Config struct {
		Bin        string `ini:"bin"`
		Remote     string `ini:"remote"`
		BaseDir    string `ini:"base_dir"`
		NTransfers int8   `ini:"n_transfers"`
		Debug      bool   `ini:"debug"`
	}

	// Dynamic variables
	DryRun bool
}

func NewRcloneWrapper(configPath string, dryRun bool) *RcloneWrapper {
	// Load config and map to struct
	rc := new(RcloneWrapper)
	cfg, err := ini.Load(configPath)
	checkError(err)

	err = cfg.Section("rclone").MapTo(&rc.Config)
	checkError(err)

	// Sanity check

	rc.DryRun = dryRun

	return rc
}

// Execute arbitrary command with rclone and output to stdout
func (rc *RcloneWrapper) exec(args ...string) {
	args = rc.prepareArgs(args)

	cmd := &exec.Cmd{
		Path:   rc.Config.Bin,
		Args:   args,
		Stdout: os.Stdout,
		Stdin:  os.Stdout,
		Stderr: os.Stderr,
	}

	DebugLogger.Printf("Executing `%s`", cmd.String())
	err := cmd.Run()
	checkError(err)
}

func (rc *RcloneWrapper) prepareArgs(args []string) []string {
	// Bin always has to be first argument
	baseArgs := []string{
		rc.Config.Bin,
		"--transfers", strconv.Itoa(int(rc.Config.NTransfers)),
		"--progress",
	}

	if rc.Config.Debug {
		baseArgs = append(baseArgs, "--verbose")
	}

	if rc.DryRun {
		baseArgs = append(baseArgs, "--dry-run")
	}

	return append(baseArgs, args...)
}

func (rc *RcloneWrapper) preparePath(path string) string {
	// Remove trailing slash if necessary
	baseDir := removeTrailSlash(rc.Config.BaseDir)
	return fmt.Sprintf("%s:%s/%s", rc.Config.Remote, baseDir, path)
}

func (rc *RcloneWrapper) List() {
	rc.exec("tree", "--dirs-only", rc.preparePath(""))
}
