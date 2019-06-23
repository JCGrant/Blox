package main

import (
	"log"

	"github.com/JCGrant/Blox/router"
	"github.com/JCGrant/Blox/wrapper"
)

func main() {
	w, err := wrapper.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Stdout.Close()
	defer w.Stdin.Close()

	go router.ParseCmdOutput(w.Stdout, w.Stdin)

	err = w.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
