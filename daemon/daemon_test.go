package daemon

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/pkg/deck"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	at := assert.New(t)

	t.Run("worker", func(t *testing.T) {
		deck.SetupEnvs(deck.Envs{envDaemonWorker: ""})
		defer deck.TeardownEnvs()

		Run()
	})

	t.Run("main success", func(t *testing.T) {
		config.Set("daemon.tries", 1)

		deck.SetupCmd()
		defer deck.TeardownCmd()
		deck.SetupOsExit(func(code int) {
			if code == 0 {
				deck.SetupEnvs(deck.Envs{envDaemon: ""})
			} else {
				at.Equal(1, code)
			}
		})
		defer deck.TeardownOsExit()
		defer deck.TeardownEnvs()

		Run()
	})

	t.Run("main error", func(t *testing.T) {
		deck.SetupCmdError()
		defer deck.TeardownCmd()

		at.Panics(Run)
	})

	t.Run("break master", func(t *testing.T) {
		config.Set("daemon.tries", 1)

		deck.SetupCmdError()
		defer deck.TeardownCmd()
		deck.SetupOsExit()
		defer deck.TeardownOsExit()

		run()
	})
}

func TestSpawn(t *testing.T) {
	at := assert.New(t)

	t.Run("skip", func(t *testing.T) {
		deck.SetupEnvs(deck.Envs{envDaemon: ""})
		defer deck.TeardownEnvs()

		cmd, err := spawn(true)
		at.Nil(err)
		at.Nil(cmd)
	})

	t.Run("redirect output", func(t *testing.T) {
		deck.SetupEnvs(deck.Envs{envDaemon: ""})
		defer deck.TeardownEnvs()

		deck.SetupCmdError()
		defer deck.TeardownCmd()

		stdoutLogFile, stderrLogFile = os.Stdout, os.Stderr

		cmd, err := spawn(false)

		at.NotNil(err)
		at.NotNil(cmd)
	})
}

func TestSetupLogFiles(t *testing.T) {
	at := assert.New(t)

	f, err := ioutil.TempFile("", "")
	at.Nil(err)
	defer func() {
		_ = os.Remove(f.Name())
	}()

	t.Run("success", func(t *testing.T) {
		config.Set("daemon.stdoutLogFile", f.Name())
		config.Set("daemon.stderrLogFile", f.Name())

		setupLogFiles()

		at.NotNil(stdoutLogFile)
		at.NotNil(stderrLogFile)
	})

	t.Run("stdout panic", func(t *testing.T) {
		config.Set("daemon.stdoutLogFile", ".")

		at.Panics(setupLogFiles)
	})

	t.Run("stderr panic", func(t *testing.T) {
		config.Set("daemon.stdoutLogFile", f.Name())
		config.Set("daemon.stderrLogFile", ".")

		at.Panics(setupLogFiles)
	})
}

func TestTeardownLogFiles(t *testing.T) {
	f := os.NewFile(1, "")
	stdoutLogFile, stderrLogFile = f, f
	defer func() { stdoutLogFile, stderrLogFile = nil, nil }()

	teardownLogFiles()
}

func TestHelperCommand(t *testing.T) {
	deck.HandleCommand(func(args []string, expectStderr bool) {})
}
