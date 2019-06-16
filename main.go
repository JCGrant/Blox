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

	"github.com/JCGrant/Blox/plotter"
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

		command, args, isCommand := parseLine(line)
		if isCommand {
			handleCommand(command, args, w)
		}
	}
}

func parseLine(line string) (string, []string, bool) {
	tokens := strings.Split(line, " ")
	if len(tokens) < 5 {
		return "", nil, false
	}
	command := tokens[4]
	if !strings.HasPrefix(command, commandChar) {
		return "", nil, false
	}
	command = strings.TrimPrefix(command, commandChar)
	var args []string
	if len(tokens) > 5 {
		args = tokens[5:]
	}
	return command, args, true
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
	case "plot":
		handlePlotCommand(args, w)
	}
}

func handleTimeCommand(args []string, w io.Writer) {
	timeStr := "0"
	if len(args) > 0 {
		timeStr = args[0]
	}
	fmt.Fprintf(w, "say setting time to %s\n", timeStr)
	fmt.Fprintf(w, "time set %s\n", timeStr)
}

func handlePlotCommand(args []string, w io.Writer) {
	block := args[0]
	oldBlockHandling := "replace"
	coords := plotter.Parse(strings.Join(args[1:], " "))
	fmt.Fprintf(w, "say plotting function\n")
	for _, c := range coords {
		x, y, z := int(c.X), int(c.Y), int(c.Z)
		fmt.Fprintf(w, "setblock %d %d %d minecraft:%s %s\n", x, y, z, block, oldBlockHandling)
	}
}
