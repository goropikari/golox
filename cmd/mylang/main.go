package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/goropikari/mylang"
)

func main() {
	runtime := mylang.NewRuntime()

	if len(os.Args) > 2 {
		fmt.Println("Usage: mylang [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1], runtime)
	} else {
		runPrompt(runtime)
	}
}

func runFile(filepath string, r *mylang.Runtime) {
	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(source)
	r.Run(bytes.NewBuffer(source))

	if r.HadError {
		os.Exit(65)
	}
}

func runPrompt(r *mylang.Runtime) {
	stdin := bufio.NewReader(os.Stdin)
	buf := &bytes.Buffer{}

	for {
		inBlock := false
		for {
			fmt.Print(prompt(inBlock))
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
			buf.Write([]byte{'\n'})

			if !canContinueRead(string(line)) {
				break
			}

			inBlock = true
		}

		r.Run(buf)
		r.HadError = false
		// fmt.Println(buf.Bytes())
		// fmt.Print(buf.String())
		buf.Reset()
	}
}

func canContinueRead(line string) bool {
	if unicode.IsSpace(int32(line[0])) {
		return true
	}

	blocks := []string{"def", "class", "if", "elseif", "else", "fun"}
	for _, v := range blocks {
		if strings.HasPrefix(line, v) {
			return true
		}
	}
	return false
}

func prompt(inBlock bool) string {
	if inBlock {
		return "... "
	}
	return ">>> "
}
