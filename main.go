package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var minecraftServerDir = "mc-server"
var minecraftServerRunParams = "-Xmx1024M -Xms1024M -jar server.jar nogui"
var commandChar = "!"

type minecraftLogger struct {
}

func (l *minecraftLogger) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%s] [Blox/INFO]: %s", time.Now().UTC().Format("15:04:05"), string(bytes))
}

var logger = log.New(&minecraftLogger{}, "", 0)

func main() {
	cmd := exec.Command("java", strings.Split(minecraftServerRunParams, " ")...)
	cmd.Dir = minecraftServerDir
	cmd.Stderr = os.Stderr

	cmdOutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal(err)
	}

	cmdInPipe, err := cmd.StdinPipe()
	if err != nil {
		logger.Fatal(err)
	}
	defer cmdInPipe.Close()

	go parseCmdOutput(cmdOutPipe, cmdInPipe)

	go io.Copy(cmdInPipe, os.Stdin)

	err = cmd.Start()
	if err != nil {
		logger.Fatal(err)
	}

	cmd.Wait()
	if err != nil {
		logger.Fatal(err)
	}
}

func parseCmdOutput(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			logger.Print(err)
		}
		line := string(lineBytes)
		fmt.Println(line)

		commandTokens := strings.Split(line, " ")[4:]
		command, args := commandTokens[0], commandTokens[1:]
		if isCommand(command) {
			command = strings.TrimPrefix(command, commandChar)
			handleCommand(command, args, w)
		}
	}
}

func isCommand(command string) bool {
	return strings.HasPrefix(command, commandChar)
}

func handleCommand(command string, args []string, w io.Writer) {
	logger.Print("handling custom command")
	defer func() {
		if r := recover(); r != nil {
			logger.Print(r)
		}
	}()
	switch command {
	case "time":
		handleTimeCommand(args, w)
	}
}

func handleTimeCommand(args []string, w io.Writer) {
	timeStr := "0"
	if len(args) > 0 {
		timeStr = args[1]
	}
	fmt.Fprintf(w, "say setting time to %s\n", timeStr)
	fmt.Fprintf(w, "time set %s\n", timeStr)
}
