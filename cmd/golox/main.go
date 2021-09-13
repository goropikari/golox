package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/goropikari/golox"
)

func main() {
	runtime := golox.NewRuntime()

	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1], runtime)
	} else {
		runPrompt(runtime)
	}
}

func runFile(path string, r *golox.Runtime) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	r.BasePath = filepath.Dir(path)

	// fmt.Println(source)
	r.Run(bytes.NewBuffer(source))

	if r.HadError {
		os.Exit(65)
	}
	if r.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt(r *golox.Runtime) {
	stdin := bufio.NewReader(os.Stdin)
	buf := &bytes.Buffer{}

	for {
		fmt.Print(">>> ")
		line, _, err := stdin.ReadLine()
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			log.Fatal(err)
		}

		if len(line) == 0 {
			break
		}

		buf.Write(line)

		r.Run(buf)
		r.HadError = false
		buf.Reset()
	}
}

func prompt(inBlock bool) string {
	if inBlock {
		return "... "
	}
	return ">>> "
}
