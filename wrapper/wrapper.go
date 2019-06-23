package wrapper

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var minecraftServerDir = "mc-server"
var minecraftServerRunParams = "-Xmx1024M -Xms1024M -jar server.jar nogui"

type minecraftLogger struct {
}

func (l *minecraftLogger) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%s] [Blox/INFO]: %s", time.Now().UTC().Format("15:04:05"), string(bytes))
}

// Wrapper wraps the Minecraft Server process
type Wrapper struct {
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	cmd    *exec.Cmd
	logger *log.Logger
}

// New sets up the Minecraft Server process
func New() (*Wrapper, error) {
	cmd := exec.Command("java", strings.Split(minecraftServerRunParams, " ")...)
	cmd.Dir = minecraftServerDir
	cmd.Stderr = os.Stderr
	cmdOutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "getting stdout pipe failed")
	}
	cmdInPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "getting stdin pipe failed")
	}
	return &Wrapper{
		Stdout: cmdOutPipe,
		Stdin:  cmdInPipe,
		cmd:    cmd,
		logger: log.New(&minecraftLogger{}, "", 0),
	}, nil
}

// Start starts the Minecraft Server process
func (w *Wrapper) Start() error {
	// Allow user to write to Minecraft Process via Wrapper
	go io.Copy(w.Stdin, os.Stdin)
	err := w.cmd.Start()
	if err != nil {
		return errors.Wrap(err, "starting Minecraft failed")
	}
	err = w.cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "waiting for Minecraft to exit failed")
	}
	return nil
}

// Printf printfs
func (w *Wrapper) Printf(format string, v ...interface{}) {
	w.logger.Printf(format, v...)
}

// Print prints
func (w *Wrapper) Print(v ...interface{}) {
	w.logger.Print(v...)
}

// Println printlns
func (w *Wrapper) Println(v ...interface{}) {
	w.logger.Println(v...)
}
