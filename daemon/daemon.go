package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/pkg/deck"
)

const envDaemon = "DAWN_DAEMON"
const envDaemonWorker = "DAWN_DAEMON_WORKER"

var stdoutLogFile *os.File
var stderrLogFile *os.File
var normalRunningTime = time.Second * 10
var osExit = deck.OsExit
var execCommand = deck.ExecCommand

func Run() {
	if isWorker() {
		return
	}

	// Panic if the initial spawned daemon process has error
	if _, err := spawn(true); err != nil {
		panic(fmt.Sprintf("dawn: failed to run in daemon mode: %s", err))
	}

	setupLogFiles()
	defer teardownLogFiles()

	run()
}

func run() {
	var (
		cmd    *exec.Cmd
		err    error
		count  int
		start  time.Time
		max    = config.GetInt("daemon.tries", 10)
		logger = log.New(stderrLogFile, "", log.LstdFlags)
	)

	for {
		if count++; count > max {
			break
		}

		start = time.Now()
		if cmd, err = spawn(false); err != nil {
			continue
		}

		err = cmd.Wait()

		logger.Printf("dawn: (pid:%d)%v exist with err: %v", cmd.Process.Pid, cmd.Args, err)

		if time.Since(start) > normalRunningTime {
			// reset count
			count = 0
		}
	}

	logger.Printf("dawn: already attempted %d times", max)

	osExit(1)
}

func spawn(skip bool) (cmd *exec.Cmd, err error) {
	if isDaemon() && skip {
		return
	}

	args, env := setupArgsAndEnv()

	cmd = execCommand(args[0], args[1:]...)
	cmd.Env = append(cmd.Env, env...)
	cmd.SysProcAttr = newSysProcAttr()

	if isDaemon() {
		if stdoutLogFile != nil {
			cmd.Stdout = stdoutLogFile
		}

		if stderrLogFile != nil {
			cmd.Stderr = stderrLogFile
		}
	}

	if err = cmd.Start(); err != nil {
		return
	}

	// Exit main process
	if !isDaemon() {
		osExit(0)
	}

	return
}

func setupArgsAndEnv() ([]string, []string) {
	args, env := os.Args, os.Environ()
	if !isDaemon() {
		args = append(args, "master process dawn")
		env = append(env, envDaemon+"=")
	} else if !isWorker() {
		args[len(args)-1] = "worker process"
		env = append(env, envDaemonWorker+"=")
	}

	return args, env
}

func isDaemon() bool {
	_, ok := os.LookupEnv(envDaemon)
	return ok
}

func isWorker() bool {
	_, ok := os.LookupEnv(envDaemonWorker)
	return ok
}

func setupLogFiles() {
	var err error
	if f := config.GetString("daemon.stdoutLogFile"); f != "" {
		if stdoutLogFile, err = os.OpenFile(filepath.Clean(f), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			panic(fmt.Sprintf("dawn: failed to open stdout log file %s: %s", f, err))
		}
	}

	if f := config.GetString("daemon.stderrLogFile"); f != "" {
		if stderrLogFile, err = os.OpenFile(filepath.Clean(f), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			panic(fmt.Sprintf("dawn: failed to open stderr log file %s: %s", f, err))
		}
	}
}

func teardownLogFiles() {
	if stdoutLogFile != nil {
		_ = stdoutLogFile.Close()
	}

	if stderrLogFile != nil {
		_ = stderrLogFile.Close()
	}
}
