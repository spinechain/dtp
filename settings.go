package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lithammer/shortuuid/v3"
)

type SpineSettings struct {
	ServerPort uint
	ClientID   string
	DataFolder string
	ShowUI     bool
}

var AppSettings SpineSettings

func LoadDefaultSettings() {
	AppSettings.ServerPort = 9100

	AppSettings.ClientID = shortuuid.New()
	AppSettings.ShowUI = true

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	AppSettings.DataFolder = filepath.Join(filepath.Dir(ex), "data")

	if _, err := os.Stat(AppSettings.DataFolder); os.IsNotExist(err) {
		err := os.Mkdir(AppSettings.DataFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func LoadSettings() {

	LoadDefaultSettings()

	settings_file := filepath.Join(AppSettings.DataFolder, "settings.json")
	fmt.Println("Loading settings from " + settings_file)

	file, err := os.Open(settings_file)
	if err != nil {
		return
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppSettings)
	if err != nil {
		fmt.Println("error:", err)
	}

}

func SaveSettings() {

	file, err := os.Create(filepath.Join(AppSettings.DataFolder, "settings.json"))
	if err != nil {
		fmt.Println("Error while writing:", err)
		return
	}

	jdata1, _ := json.MarshalIndent(AppSettings, "", " ")

	file.Write(jdata1)

	defer file.Close()
}
