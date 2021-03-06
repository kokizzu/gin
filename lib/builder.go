package gin

import (
	"os/exec"
	"runtime"
	"strings"
	"fmt"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
	SetErrors(str string)
}

type builder struct {
	dir      string
	binary   string
	errors   string
	useGodep bool
}

func NewBuilder(dir string, bin string, useGodep bool) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") {
			// check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, useGodep: useGodep}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	var command *exec.Cmd
	if b.useGodep {
		command = exec.Command("godep", "go", "build", "-o", b.binary)
	} else {
		command = exec.Command("go", "build", "-i", `-x`, `-toolexec`, "/usr/bin/time -f '\t%e s\t%M KB'", "-o", b.binary)
	}
	command.Dir = b.dir
	res, _ := command.CombinedOutput()
	out := string(res)
	fmt.Print(out)

	if !command.ProcessState.Success() {
		b.errors = out
		return fmt.Errorf(`%s`, out)
	} else {
		b.errors = ``
		return nil
	}
}

func (b *builder) SetErrors(str string) {
	b.errors = str
}