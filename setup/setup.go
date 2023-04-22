package setup

import (
	"fmt"
	"os"
)

var HOME_DIR string = os.Getenv("HOME")
var GDB_INIT_FILE string = HOME_DIR + "/.gdbinit"

func Setgdb() bool {
	err := os.Truncate(HOME_DIR+"/.gdbinit", 0)
	if !(err == nil) {
		fmt.Println("~/.gdbinit is not found")
		return false
	} else {
		return true
	}
}

func Setpeda(dir string) bool {
	err := os.Truncate(GDB_INIT_FILE, 0)
	if !(err == nil) {
		fmt.Fprintln(os.Stderr, "[x] ~/.gdbinit is not found")
		return false
	} else {
		var command string = "source " + dir + "/peda.py"
		var byted_command []byte = []byte(command)
		os.WriteFile(GDB_INIT_FILE, byted_command, 0644)
		return true
	}
}

func Setgef(dir string) bool {
	err := os.Truncate(GDB_INIT_FILE, 0)
	if !(err == nil) {
		fmt.Fprintln(os.Stderr, "[x] ~/.gdbinit is not found")
		return false
	} else {
		var command string = "source " + dir + "/.gdbinit-gef.py"
		var byted_command []byte = []byte(command)
		os.WriteFile(GDB_INIT_FILE, byted_command, 0644)
		return true
	}
}

func Setpwndbg(dir string) bool {
	err := os.Truncate(GDB_INIT_FILE, 0)
	if !(err == nil) {
		fmt.Fprintln(os.Stderr, "[x] ~/.gdbinit is not found")
		return false
	} else {
		var command1 string = "set debuginfod enabled on"
		var command2 string = "source " + dir + "/gdbinit.py"
		var byted_command1 []byte = []byte(command1)
		var byted_command2 []byte = []byte(command2)
		os.WriteFile(GDB_INIT_FILE, byted_command1, 0644)
		os.WriteFile(GDB_INIT_FILE, byted_command2, 0644)
		return true
	}
}
