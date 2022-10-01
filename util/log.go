package util

import (
	"log"
	"os"
	"path/filepath"
)

var file *os.File

func CreateLog(dataDir string) {
	// create log file
	var err error
	file, err = os.OpenFile(filepath.Join(dataDir, "spine.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Set log out put and enjoy :)
	log.SetOutput(file)

	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func CloseLogFile() {
	log.SetOutput(os.Stdout)
	file.Close()
}
