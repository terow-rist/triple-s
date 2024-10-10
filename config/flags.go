package config

import (
	"flag"
	"fmt"
)

var (
	PortNumber string
	Directory  string
)

func init() {
	flag.StringVar(&PortNumber, "port", "8080", "Port number")
	flag.StringVar(&Directory, "dir", "", "Path to the directory")

	helpMessage :=
		`Simple Storage Service.

**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`

	flag.Usage = func() {
		fmt.Println(helpMessage)
	}
	flag.Parse()
}
