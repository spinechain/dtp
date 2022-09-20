package tasktypes

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"spinedtp/util"
	"strings"
	"time"

	_ "embed"
)

//go:embed ping.bat
var ping_bat string

//go:embed ping.sh
var ping_sh string

//go:embed stable_diffusion.sh
var sd_sh string

var isRunning bool = false

// Struture for tasks to be executed
type TaskToExecute struct {
	TaskType    string
	TaskID      string
	DataFolder  string
	Prompt      string
	Complete    bool
	ResultFiles []string
	MimeType    string
	ResultError error
	Sent        bool
}

type TaskType struct {
	name              string
	trigger           string
	scriptFolder      string
	execFolder        string
	windowsScriptName string
	linuxScriptName   string
	windowsFullScript string
	linuxFullScript   string
	macFullScript     string
	macScriptName     string
	scriptPath        string
	fullScript        string
	outputExtension   string
	outputSubpath     string
	ValidationRegex   string
	RegexCompiled     *regexp.Regexp
}

var TaskTypes []TaskType

// List with tasks to execute
var TasksToExecute []TaskToExecute
var TaskForSubmissionAvailable chan int

func Init(DataFolder string) error {
	// copy the emebeded scripts to the data folder

	// Create stable diffusion task type
	var sd TaskType
	sd.name = "sd"
	sd.execFolder = "/home/mark/stable-diffusion"
	sd.scriptFolder = filepath.Join(DataFolder, "scripts")
	sd.windowsScriptName = "stable_diffusion.bat"
	sd.linuxFullScript = sd_sh
	sd.linuxScriptName = "stable_diffusion.sh"
	sd.macScriptName = "stable_diffusion.sh"
	sd.outputSubpath = "samples"
	sd.outputExtension = ".png;.txt"
	sd.trigger = "draw"
	TaskTypes = append(TaskTypes, sd)

	// Create ping task type
	var pg TaskType
	pg.name = "ping"
	pg.execFolder = ""
	pg.scriptFolder = filepath.Join(DataFolder, "scripts")
	pg.windowsScriptName = "ping.bat"
	pg.linuxFullScript = ping_sh
	pg.windowsFullScript = ping_bat
	pg.linuxScriptName = "ping.sh"
	pg.macScriptName = "ping.sh"
	pg.outputSubpath = ""
	pg.trigger = "ping"
	pg.outputExtension = ".txt"
	pg.ValidationRegex = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
	pg.RegexCompiled, _ = regexp.Compile(pg.ValidationRegex)

	TaskTypes = append(TaskTypes, pg)

	// Loop over all task types
	for i, taskType := range TaskTypes {

		ops := runtime.GOOS
		switch ops {
		case "windows":
			taskType.scriptPath = filepath.Join(taskType.scriptFolder, taskType.windowsScriptName)
			taskType.fullScript = taskType.windowsFullScript
		case "darwin":
			taskType.scriptPath = filepath.Join(taskType.scriptFolder, taskType.macScriptName)
			taskType.fullScript = taskType.macFullScript
		case "linux":
			taskType.scriptPath = filepath.Join(taskType.scriptFolder, taskType.linuxScriptName)
			taskType.fullScript = taskType.linuxFullScript
		default:
			continue
		}

		TaskTypes[i] = taskType

		// create the scripts folder if it does not exist
		if _, err := os.Stat(taskType.scriptFolder); os.IsNotExist(err) {
			err := os.Mkdir(taskType.scriptFolder, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		// delete existing scripts
		if _, err := os.Stat(taskType.scriptPath); !os.IsNotExist(err) {
			err := os.Remove(taskType.scriptPath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// check if the files exist
		if _, err := os.Stat(taskType.scriptPath); os.IsNotExist(err) {
			// file does not exist
			err := ioutil.WriteFile(taskType.scriptPath, []byte(taskType.fullScript), 0644)
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	return nil
}

func AddToTaskExecutionQueue(dataFolder string, taskType string, taskID string, prompt string) error {
	util.PrintYellow("Adding to execution queue: " + taskType)

	// add to the queue
	TasksToExecute = append(TasksToExecute, TaskToExecute{TaskType: taskType, DataFolder: dataFolder, Prompt: prompt, TaskID: taskID, Complete: false})

	// add to queue
	if !isRunning {
		go RunTaskExecutionProcess()
	} else {
		util.PrintYellow("Task exec script is already running")
	}

	return nil
}

func CompleteTask(task *TaskToExecute, resultFiles []string, resultError error) {

	// print
	util.PrintYellow("Task complete: " + task.TaskType)

	task.Complete = true
	task.ResultFiles = resultFiles
	task.ResultError = resultError

	TaskForSubmissionAvailable <- 1
}

// Get next open task
func GetNextTask() *TaskToExecute {
	for i, task := range TasksToExecute {
		if !task.Complete {
			return &TasksToExecute[i]
		}
	}

	return nil
}

func RunTaskExecutionProcess() error {

	isRunning = true
	defer func() { isRunning = false }()

	util.PrintYellow("Starting task execution process")

	// big loop
	for te := GetNextTask(); te != nil; te = GetNextTask() {

		outputDir := filepath.Join(te.DataFolder, "output")
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			err := os.Mkdir(outputDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		// find the task type
		var taskType *TaskType
		for _, tt := range TaskTypes {
			if tt.name == te.TaskType {
				taskType = &tt
			}
		}

		if taskType == nil {
			CompleteTask(&TasksToExecute[0], nil, errors.New("task type not supported"))
			continue
		}

		te.Prompt = strings.Replace(te.Prompt, taskType.trigger, "", 1)
		te.Prompt = strings.TrimSpace(te.Prompt)

		if !taskType.RegexCompiled.MatchString(te.Prompt) {
			util.PrintRed("Cannot match string")
			continue
		}

		var cmd *exec.Cmd
		if filepath.Ext(taskType.scriptPath) == ".bat" {
			// run the batch file
			cmd = exec.Command("cmd.exe", "/C", taskType.scriptPath, taskType.execFolder, te.Prompt, outputDir)
		} else {
			// run the shell script
			cmd = exec.Command("bash", taskType.scriptPath, taskType.execFolder, te.Prompt, outputDir)
		}

		startTime := time.Now()
		util.PrintGreen("Running task: " + te.TaskType + " at time " + startTime.Format("2006-01-02 15:04:05"))

		data, err := cmd.CombinedOutput()
		if err != nil {
			// Complete task
			CompleteTask(&TasksToExecute[0], nil, err)
			continue
		}

		if err != nil {
			// Complete task
			CompleteTask(&TasksToExecute[0], nil, err)
			continue
		}

		util.PrintGreen("Task complete: " + te.TaskType + " duration " + time.Since(startTime).String())

		// print the output
		util.PrintGreen(string(data))

		resultFiles, mimeType := FigureOutResultFile(filepath.Join(outputDir, taskType.outputSubpath), taskType.outputExtension, startTime)

		if resultFiles == nil || len(resultFiles) == 0 {
			err = errors.New("no result files found")
		}

		TasksToExecute[0].MimeType = mimeType

		// Complete task
		CompleteTask(&TasksToExecute[0], resultFiles, err)
	}

	return nil
}

// This will tell us what the result file is. If there is more than one potential result file,
// it will fail.
func FigureOutResultFile(filesFolder string, ext string, execStartTime time.Time) ([]string, string) {

	// find the result file
	files, err := ioutil.ReadDir(filesFolder)
	if err != nil {
		return nil, ""
	}

	var resultFiles []string

	var mimeType string
	// find the result file
	for _, f := range files {

		// check if the file is the right extension
		if f.ModTime().After(execStartTime) {

			extensions := strings.Split(ext, ";")

			for _, extension := range extensions {
				if filepath.Ext(f.Name()) == extension {
					resultFiles = append(resultFiles, filepath.Join(filesFolder, f.Name()))
					mimeType = extension
				}
			}
		}
	}

	switch mimeType {
	case ".png":
		mimeType = "image/png"
	case ".jpg":
		mimeType = "image/jpeg"
	case ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".mp4":
		mimeType = "video/mp4"
	case ".mp3":
		mimeType = "audio/mpeg"
	case ".wav":
		mimeType = "audio/wav"
	case ".ogg":
		mimeType = "audio/ogg"
	case ".txt":
		mimeType = "text/plain"
	case ".html":
		mimeType = "text/html"
	case ".json":
		mimeType = "application/json"
	case ".xml":
		mimeType = "application/xml"
	case ".zip":
		mimeType = "application/zip"
	case ".pdf":
		mimeType = "application/pdf"
	case ".doc":
		mimeType = "application/msword"
	}

	return resultFiles, mimeType
}
