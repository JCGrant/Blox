package main

import (
	"log"

	"github.com/JCGrant/Blox/proxy"
)

func main() {
	p := proxy.New()
	err := p.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
