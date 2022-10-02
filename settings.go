package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	tasknet "spinedtp/tasknet"
	"spinedtp/tasktypes"
	"spinedtp/ui"
	"spinedtp/util"

	"github.com/kirsle/configdir"
	"github.com/lithammer/shortuuid/v3"
)

type SpineSettings struct {
	ServerPort      uint
	ListenAddress   string
	ClientID        string
	ConfigFolder    string
	TaskTypesFolder string
	LogFolder       string
	DbFolder        string
	ShowUI          bool
	RouteOnly       bool // In this case it acts as a router only, no tasks are executed
}

var AppSettings SpineSettings

func LoadDefaultSettings() {
	AppSettings.ServerPort = 59143
	AppSettings.ListenAddress = "127.0.0.1"
	AppSettings.ClientID = shortuuid.New()
	AppSettings.ShowUI = true
	AppSettings.RouteOnly = false

	configPath := configdir.LocalConfig("spinechain")
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		panic(err)
	}

	AppSettings.ConfigFolder = configPath

	if _, err := os.Stat(AppSettings.ConfigFolder); os.IsNotExist(err) {
		err := os.Mkdir(AppSettings.ConfigFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	AppSettings.TaskTypesFolder = filepath.Join(AppSettings.ConfigFolder, "tasktypes")

	if _, err := os.Stat(AppSettings.TaskTypesFolder); os.IsNotExist(err) {
		err := os.Mkdir(AppSettings.TaskTypesFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	appConfigDir := configdir.LocalCache("spinechain")
	AppSettings.LogFolder = filepath.Join(appConfigDir, "logs")
	AppSettings.DbFolder = filepath.Join(appConfigDir, "db")

	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(appConfigDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(AppSettings.LogFolder); os.IsNotExist(err) {
		err := os.Mkdir(AppSettings.LogFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(AppSettings.DbFolder); os.IsNotExist(err) {
		err := os.Mkdir(AppSettings.DbFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	tasktypes.Init(AppSettings.TaskTypesFolder)
	tasknet.TasksToExecute = &tasktypes.TasksToExecute
}

func LoadSettings() string {

	LoadDefaultSettings()

	settings_file := filepath.Join(AppSettings.ConfigFolder, "settings.json")
	fmt.Println("Loading settings from " + settings_file)

	file, err := os.Open(settings_file)
	if err != nil {
		return settings_file
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppSettings)
	if err != nil {
		fmt.Println("error:", err)
	}

	util.CreateLog(AppSettings.LogFolder)

	return settings_file
}

func SaveSettings() {

	file, err := os.Create(filepath.Join(AppSettings.ConfigFolder, "settings.json"))
	if err != nil {
		fmt.Println("Error while writing:", err)
		return
	}

	jdata1, _ := json.MarshalIndent(AppSettings, "", " ")

	file.Write(jdata1)

	defer file.Close()
}

func SetNetworkSettings() {

	tasknet.NetworkSettings.ServerHost = AppSettings.ListenAddress
	tasknet.NetworkSettings.MyPeerID = AppSettings.ClientID
	tasknet.NetworkSettings.ServerPort = AppSettings.ServerPort
	tasknet.NetworkSettings.OnStatusUpdate = Event_StatusUpdate
	tasknet.NetworkSettings.BidTimeout = 5
	tasknet.NetworkSettings.AcceptedBidsPerTask = 3
	tasknet.NetworkSettings.TaskTypeFolder = AppSettings.TaskTypesFolder
	tasknet.NetworkSettings.DbFolder = AppSettings.DbFolder
	tasknet.NetworkSettings.RouteOnly = AppSettings.RouteOnly

	tasknet.NetworkCallbacks.OnTaskReceived = nil // s.OnNewTaskReceived
	tasknet.NetworkCallbacks.OnTaskApproved = nil //s.OnNetworkTaskApproval
	tasknet.NetworkCallbacks.OnTaskResult = Event_TaskResultReceived

	tasknet.LoadDefaultPeerTable(default_peers)

	UpdateInfoStatusBar()
}

func UpdateInfoStatusBar() {
	ui.UpdateStatusBar("🤖 "+AppSettings.ClientID, 2)
}
