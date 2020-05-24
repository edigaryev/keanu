package main

import (
	"flag"
	"github.com/edigaryev/keanu/preprocessor"
	"log"
	"os"
)

const (
	ArgInputPath = iota
	ArgOutputPath
	ArgMax
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || len(flag.Args()) > ArgMax {
		log.Fatal("usage: keanu input.yaml [output.yaml]")
	}

	inputPath := flag.Arg(ArgInputPath)
	p, err := preprocessor.NewFromFile(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Determine the output destination: stdout or file
	var output *os.File
	if len(flag.Args()) == ArgMax {
		outputPath := flag.Arg(ArgOutputPath)
		output, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	out, err := p.Dump()
	if err != nil {
		log.Fatal(err)
	}

	_, err = output.Write(out)
	if err != nil {
		log.Fatal(err)
	}

	err = output.Close()
	if err != nil {
		log.Fatal(err)
	}
}
