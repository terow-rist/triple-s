package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	PortNumber string
	Directory  string
)

func init() {
	flag.StringVar(&PortNumber, "port", "8080", "Port number")
	flag.StringVar(&Directory, "dir", "data", "Path to the directory")

	helpMessage := `Simple Storage Service.

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

func ValidateDirectory() error {
	// checking that --dir=path exists
	if _, err := os.Stat(Directory); os.IsNotExist(err) {
		err = os.Mkdir(Directory, 0o755)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
	}
	// checking that '--dir=' is standard or not
	if isStandardPackage(Directory) {
		return errors.New("Error: directory(--dir=) cannot be one of the used ones.")
	}
	return nil
}

func isStandardPackage(packageName string) bool {
	return packageName == "cmd" || packageName == "config" || packageName == "internal"
}
