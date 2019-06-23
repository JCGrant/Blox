package plugins

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/JCGrant/Blox/plotter"
)

var exampleUsage = "!plot block x-expression, y-expression, z-expression | i <- 1..10"

func HandleCommand(command string, args []string, w io.Writer) {
	log.Print("handling custom command")
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
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
	if len(args) < 2 {
		fmt.Fprintf(w, "say usage: %s\n", exampleUsage)
		return
	}
	block := args[0]
	oldBlockHandling := "replace"
	coords, err := plotter.Parse(strings.Join(args[1:], " "))
	if err != nil {
		fmt.Fprintf(w, "say parsing plot failed: %s\n", err)
		return
	}
	fmt.Fprintf(w, "say plotting function\n")
	for _, c := range coords {
		x, y, z := int(c.X), int(c.Y), int(c.Z)
		fmt.Fprintf(w, "setblock %d %d %d minecraft:%s %s\n", x, y, z, block, oldBlockHandling)
	}
}
