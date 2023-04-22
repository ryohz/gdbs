package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cavaliergopher/grab/v3"
	"github.com/go-git/go-git/v5"
)

var HOME_DIR string = os.Getenv("HOME")
var GDB_INIT_FILE string = HOME_DIR + "/.gdbinit"

func Installpeda(dir string) bool {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/longld/peda.git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprint(os.Stderr, "[x] ")
		fmt.Fprintln(os.Stderr, err)
		return false
	} else {
		return true
	}
}

func Installgef(dir string) bool {
	_, err := grab.Get(dir+"/.gdbinit-gef.py", "https://gef.blah.cat/py")
	if err != nil {
		fmt.Fprint(os.Stderr, "[x] ")
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return true
}

func Installpwndbg(dir string) bool {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/pwndbg/pwndbg.git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprint(os.Stderr, "[x] ")
		fmt.Fprintln(os.Stderr, err)
		return false
	} else {
		err := os.Truncate(GDB_INIT_FILE, 0)
		if !(err == nil) {
			fmt.Fprintln(os.Stderr, "[x] ~/.gdbinit is not found")
			return false
		} else {

			os.Chdir(dir)
			cmd := exec.Command("./setup.sh")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Run()
			os.Chdir("../")
			return true
		}
	}
}
