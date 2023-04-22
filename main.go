package main

import (
	"fmt"
	"io/fs"
	"os"

	"ryozk/gdbs/install"

	"github.com/fatih/color"

	"ryozk/gdbs/setup"
)

func include(list []string, target string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == target {
			return true
		}
	}
	return false
}

type Gdbmod struct {
	Name                 string
	InstallDir           string
	AppDir               string
	GdbModsContainsNames []string
	StateFile            string
}

func Install(mod Gdbmod) bool {
	if mod.Name == "gdb" {
		return true
	}
	if !include(mod.GdbModsContainsNames, mod.Name) {
		var err_flag bool = false
		switch mod.Name {
		case "peda":
			if !install.Installpeda(mod.InstallDir) {
				err_flag = !err_flag
			}
		case "gef":
			if !install.Installgef(mod.InstallDir) {
				err_flag = !err_flag
			}
		case "pwndbg":
			if !install.Installpwndbg(mod.InstallDir) {
				err_flag = !err_flag
			}
		}
		if err_flag {
			fmt.Fprintln(os.Stderr, "failed to install "+mod.Name)
			return false
		} else {
			fmt.Println("installed " + mod.Name + " successfully")
			return true
		}
	} else {
		return true
	}
}

func Setup(mod Gdbmod) {
	var err_flag bool = false
	switch mod.Name {
	case "gdb":
		if !setup.Setgdb() {
			err_flag = !err_flag
		}
	case "peda":
		if !setup.Setpeda(mod.InstallDir) {
			err_flag = !err_flag
		}
	case "gef":
		if !setup.Setgef(mod.InstallDir) {
			err_flag = !err_flag
		}
	case "pwndbg":
		if !setup.Setpwndbg(mod.InstallDir) {
			err_flag = !err_flag
		}
	}
	if err_flag {
		fmt.Fprintln(os.Stderr, "[x] failed to change tool to "+mod.Name)
		return
	} else {
		fmt.Println("[*] changed tool to " + mod.Name + " successfully")
		if os.Truncate(mod.StateFile, 0) != nil {
			os.Create(mod.StateFile)
		}
		os.WriteFile(mod.AppDir+"/state.txt", []byte(mod.Name), 0644)
		return
	}
}

func main() {
	var argslen int = len(os.Args) - 1
	if argslen == 0 {
		fmt.Println("Usage:")
		fmt.Println("  command [option]")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  show			show current tool")
		fmt.Println("  gdb			set tool to legacy gdb")
		fmt.Println("  peda			set tool to gdb-peda")
		fmt.Println("  gef			set tool to gef")
		fmt.Println("  pwndbg		set tool to pwndbg")
		fmt.Println("  clean			clean application")
		return
	}
	if argslen == 1 {
		//get user's home directory
		hd, _ := os.UserHomeDir()
		var home_dir string = hd
		//set directory application generates directories and files
		//if the directory is not existing, make it
		var app_dir string = home_dir + "/.gdbs"
		var app_dir_contains []fs.DirEntry
		adc1, err := os.ReadDir(app_dir)
		if err != nil {
			os.Mkdir(app_dir, 0755)
			adc2, _ := os.ReadDir(app_dir)
			app_dir_contains = adc2
		} else {
			app_dir_contains = adc1
		}
		//get string list of names of directories in app directory
		//if "gdbmods" (dir) is not existing in this list, make this directory
		var app_dir_contains_names []string = []string{}
		for i := 0; i < len(app_dir_contains); i++ {
			app_dir_contains_names = append(app_dir_contains_names, app_dir_contains[i].Name())
		}
		if include(app_dir_contains_names, "gdbmods") {
			os.Mkdir(app_dir+"/gdbmods/", 0755)
		}
		var gdbmods_dir string = app_dir + "/gdbmods/"
		//get string list of names of directory in gdbmods directory to confirm
		//installed tools later
		var gdbmods_dir_contains_names []string = []string{}
		gdbmods_dir_contains, _ := os.ReadDir(app_dir + "/gdbmods/")
		for i := 0; i < len(gdbmods_dir_contains); i++ {
			gdbmods_dir_contains_names = append(gdbmods_dir_contains_names, gdbmods_dir_contains[i].Name())
		}
		//make list of all tools can be installed and used.
		var mod_list []string = []string{"gdb", "peda", "gef", "pwndbg"}
		//path of state file to manage current tool
		var state_file_path string = app_dir + "/state.txt"
		if !include(app_dir_contains_names, "state.txt") {
			new_state_file, err := os.Create(state_file_path)
			if err != nil {
				return
			}
			os.WriteFile(new_state_file.Name(), []byte("gdb"), 0644)
		}
		//processing of install or changing or both of each tool.
		for i := 0; i < len(mod_list); i++ {
			if os.Args[1] == mod_list[i] {
				var eachmod Gdbmod
				eachmod.Name = mod_list[i]
				eachmod.InstallDir = gdbmods_dir + eachmod.Name
				eachmod.AppDir = app_dir
				eachmod.GdbModsContainsNames = gdbmods_dir_contains_names
				eachmod.StateFile = state_file_path
				if !Install(eachmod) {
					return
				} else {
					Setup(eachmod)
				}
			}
		}
		switch os.Args[1] {
		//legacy gdb session
		case "show":
			state, err := os.ReadFile(state_file_path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "[x] the file to manage state of tools is not found")
				return
			} else {
				if include(mod_list[:], string(state)) {
					for i := 0; i < len(mod_list); i++ {
						if string(state) == mod_list[i] {
							fmt.Println("[" + color.RedString("*") + "] " + color.RedString(mod_list[i]))
						} else {
							fmt.Println("[-] " + mod_list[i])
						}
					}
					return
				} else {
					os.Truncate(state_file_path, 0)
					os.WriteFile(state_file_path, []byte("gdb"), 0644)
					return
				}
			}
		//cleaning session
		case "clean":
			os.RemoveAll(app_dir)
			fmt.Println("[*] uninstalled tools successfully")
			return
		}
		//help session
	} else {
		return
	}
}
