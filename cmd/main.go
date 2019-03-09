package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/induzo/fsm"
)

func main() {

	flag.Parse()

	g := fsm.NewGraph()

	// Open json file
	f, errOF := os.OpenFile("./example/graph.json", os.O_RDONLY, 0600)
	if errOF != nil {
		log.Fatalf("OpenFile: %v", errOF)
	}
	data, errRA := ioutil.ReadAll(f)
	if errRA != nil {
		log.Fatalf("OpenFile: %v", errRA)
	}

	// Parse into the graph
	if err := g.UnmarshalJSON(data); err != nil {
		log.Fatalf("UnmarshalJSON: %v", err)
	}

	fg, errOF := os.OpenFile("./example/graph.png", os.O_RDWR|os.O_CREATE, 0600)
	if errOF != nil {
		log.Fatalf("OpenFile: %v", errOF)
	}
	defer func() {
		if err := fg.Close(); err != nil {
			log.Fatalf("fg.Close: %v", err)
		}
	}()

	if err := g.GeneratePNG(fg); err != nil {
		log.Fatalf("GeneratePNG: %v", err)
	}
}
