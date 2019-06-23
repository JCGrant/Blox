package router

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/JCGrant/Blox/plugins"
)

var commandChar = "!"

func ParseCmdOutput(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Print(err)
		}
		line := string(lineBytes)
		fmt.Println(line)

		command, args, isCommand := parseLine(line)
		if isCommand {
			plugins.HandleCommand(command, args, w)
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
