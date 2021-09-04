package mylang

import (
	"bytes"
	"fmt"
	"os"
)

// Runtime is struct of Runtime
type Runtime struct {
	HadError bool
}

// NewRuntime is constructor of Runtime
func NewRuntime() *Runtime {
	return &Runtime{HadError: false}
}

// Run runs script
func (r *Runtime) Run(source *bytes.Buffer) {

	scanner := NewScanner(r, source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

// ErrorMessage prints error massage at stderr
func (r *Runtime) ErrorMessage(line int, message string) {
	r.Report(line, "", message)
}

// Report prints error masseg at stderr
func (r *Runtime) Report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line "+fmt.Sprint(line)+"] Error"+where+": "+message)
	r.HadError = true
}
