package tasktypes

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"spinedtp/util"
	"time"

	_ "embed"
)

//go:embed stable_diffusion.bat
var ld_bat string

//go:embed stable_diffusion.sh
var ld_sh string

var isRunning bool = false

// Struture for tasks to be executed
type TaskToExecute struct {
	TaskType    string
	TaskID      string
	DataFolder  string
	Prompt      string
	Complete    bool
	ResultFile  string
	ResultError error
	Sent        bool
}

type TaskType struct {
	name                string
	script_folder       string
	exec_folder         string
	windows_script_name string
	linux_script_name   string
	mac_script_name     string
	script_path         string
	full_script         string
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
	sd.exec_folder = "/home/mark/stable-diffusion"
	sd.script_folder = filepath.Join(DataFolder, "scripts")
	sd.windows_script_name = "stable_diffusion.bat"
	sd.linux_script_name = "stable_diffusion.sh"
	sd.mac_script_name = "stable_diffusion.sh"
	TaskTypes = append(TaskTypes, sd)

	// Loop over all task types
	for i, taskType := range TaskTypes {

		ops := runtime.GOOS
		switch ops {
		case "windows":
			taskType.script_path = filepath.Join(taskType.script_folder, taskType.windows_script_name)
			taskType.full_script = ld_bat
		case "darwin":
			taskType.script_path = filepath.Join(taskType.script_folder, taskType.mac_script_name)
		case "linux":
			taskType.script_path = filepath.Join(taskType.script_folder, taskType.linux_script_name)
			taskType.full_script = ld_sh
		default:
			continue
		}

		TaskTypes[i] = taskType

		// create the scripts folder if it does not exist
		if _, err := os.Stat(taskType.script_folder); os.IsNotExist(err) {
			err := os.Mkdir(taskType.script_folder, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		// delete existing scripts
		if _, err := os.Stat(taskType.script_path); !os.IsNotExist(err) {
			err := os.Remove(taskType.script_path)
			if err != nil {
				log.Fatal(err)
			}
		}

		// check if the files exist
		if _, err := os.Stat(taskType.script_path); os.IsNotExist(err) {
			// file does not exist
			err := ioutil.WriteFile(taskType.script_path, []byte(taskType.full_script), 0644)
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

func CompleteTask(task *TaskToExecute, resultFile string, resultError error) {

	// print
	util.PrintYellow("Task complete: " + task.TaskType)

	task.Complete = true
	task.ResultFile = resultFile
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
			CompleteTask(&TasksToExecute[0], "", errors.New("task type not supported"))
			continue
		}

		var cmd *exec.Cmd
		if filepath.Ext(taskType.script_path) == ".bat" {
			// run the batch file
			cmd = exec.Command("cmd.exe", "/C", taskType.script_path, taskType.exec_folder, te.Prompt, outputDir)
		} else {
			// run the shell script
			cmd = exec.Command("bash", taskType.script_path, taskType.exec_folder, te.Prompt, outputDir)
		}

		// stdout
		// stdout, err := cmd.StdoutPipe()

		startTime := time.Now()

		util.PrintGreen("Running task: " + te.TaskType + " at time " + startTime.Format("2006-01-02 15:04:05"))
		// err = cmd.Start()
		data, err := cmd.CombinedOutput()
		if err != nil {
			// Complete task
			CompleteTask(&TasksToExecute[0], "", err)
			continue
		}

		// read all
		// stdoutBytes, err := ioutil.ReadAll(stdout)

		// err = cmd.Wait()
		if err != nil {
			// Complete task
			CompleteTask(&TasksToExecute[0], "", err)
			continue
		}

		util.PrintGreen("Task complete: " + te.TaskType + " duration " + time.Since(startTime).String())

		// print the output
		util.PrintGreen(string(data))

		// result to string
		// resultString := string(result)
		//fmt.Println(resultString)

		//if strings.Contains(resultString, "Enjoy.") {
		//	util.PrintBlue("Latent diffusion completed successfully")

		//resultFile, err := FigureOutResultFile(outputDir, startTime)
		//	if err != nil {
		//		// Complete task
		//		CompleteTask(&TasksToExecute[0], "", err)
		//		continue
		//	}

		// Complete task
		//CompleteTask(&TasksToExecute[0], resultFile, nil)
		// }

		// TODO: this should be an error, if we don't know the command
		CompleteTask(&TasksToExecute[0], "", nil)
	}

	return nil
}

// This will tell us what the result file is. If there is more than one potential result file,
// it will fail.
func FigureOutResultFile(filesFolder string, execStartTime time.Time) (string, error) {

	// find the result file
	files, err := ioutil.ReadDir(filesFolder)
	if err != nil {
		log.Fatal(err)
	}

	// find the result file
	for _, f := range files {
		if f.ModTime().After(execStartTime) {
			// this is the result file
			util.PrintBlue("Result file: " + f.Name())
			return f.Name(), nil
		}
	}

	return "", errors.New("no result file found")
}
