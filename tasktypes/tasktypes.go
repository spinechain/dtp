package tasktypes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"spinedtp/util"
	"strings"
	"time"

	_ "embed"
)

//go:embed latent_diffusion.bat
var ld_bat string

//go:embed latent_diffusion.sh
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

// List with tasks to execute
var TasksToExecute []TaskToExecute
var TaskForSubmissionAvailable chan int

func CopySripts(DataFolder string) error {
	// copy the emebeded scripts to the data folder
	// latent diffusion
	// latent diffusion
	ld_bat_file := filepath.Join(DataFolder, "scripts", "latent_diffusion.bat")
	ld_sh_file := filepath.Join(DataFolder, "scripts", "latent_diffusion.sh")

	// create the scripts folder if it does not exist
	if _, err := os.Stat(filepath.Join(DataFolder, "scripts")); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(DataFolder, "scripts"), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// delete existing scripts
	if _, err := os.Stat(ld_bat_file); !os.IsNotExist(err) {
		err := os.Remove(ld_bat_file)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(ld_sh_file); !os.IsNotExist(err) {
		err := os.Remove(ld_sh_file)
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error
	// check if the files exist
	if _, err := os.Stat(ld_bat_file); os.IsNotExist(err) {
		// file does not exist
		err := ioutil.WriteFile(ld_bat_file, []byte(ld_bat), 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	if _, err := os.Stat(ld_sh_file); os.IsNotExist(err) {
		// file does not exist
		err := ioutil.WriteFile(ld_sh_file, []byte(ld_sh), 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return err
}

func AddToTaskExecutionQueue(dataFolder string, taskType string, taskID string, prompt string) error {
	util.PrintYellow("Adding to execution queue: " + taskType)

	// add to the queue
	TasksToExecute = append(TasksToExecute, TaskToExecute{TaskType: taskType, DataFolder: dataFolder, Prompt: prompt, TaskID: taskID, Complete: false})

	// add to queue
	if !isRunning {
		go RunTaskExecutionProcess()
	} else {
		util.PrintYellow("Latent diffusion is already running")
	}

	return nil
}

func CompleteTask(task *TaskToExecute, resultFile string, resultError error) {
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

		var shellScriptName string
		ops := runtime.GOOS
		switch ops {
		case "windows":
			if te.TaskType == "ld" {
				// check if linux or windows
				shellScriptName = "latent_diffusion.bat"
			}
		case "darwin":
			util.PrintRed("LD is not supported on mac")
		case "linux":
			if te.TaskType == "ld" {
				// check if linux or windows
				shellScriptName = "latent_diffusion.sh"
			}
		default:
			fmt.Printf("%s.\n", ops)
		}

		if shellScriptName == "" {
			CompleteTask(&TasksToExecute[0], "", errors.New("task type not supported"))
			continue
		}

		var cmd *exec.Cmd
		if filepath.Ext(shellScriptName) == ".bat" {
			// run the batch file
			cmd = exec.Command("cmd.exe", "/C", filepath.Join(te.DataFolder, "scripts", shellScriptName), te.Prompt, outputDir)
		} else {
			// run the shell script
			cmd = exec.Command("bash", filepath.Join(te.DataFolder, "scripts", shellScriptName), te.Prompt, outputDir)
		}

		stdout, err := cmd.StdoutPipe()

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadAll(stdout)

		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", string(data))

		/*
			startTime := time.Now()

			util.PrintGreen("Running task: " + te.TaskType + " at time " + startTime.Format("2006-01-02 15:04:05"))
			result, err := cmd.CombinedOutput()
			if err != nil {
				// Complete task
				CompleteTask(&TasksToExecute[0], "", err)
				continue
			}

			util.PrintGreen("Task complete: " + te.TaskType + " duration " + time.Since(startTime).String())
		*/
		// result to string
		resultString := string(data)

		// search for text in result
		if strings.Contains(resultString, "not a valid Win32") {
			util.PrintRed("The latent diffusion script is not a valid Win32 application")
		}

		fmt.Println(resultString)

		if strings.Contains(resultString, "Enjoy.") {
			util.PrintBlue("Latent diffusion completed successfully")

			//resultFile, err := FigureOutResultFile(outputDir, startTime)
			if err != nil {
				// Complete task
				CompleteTask(&TasksToExecute[0], "", err)
				continue
			}

			// Complete task
			//CompleteTask(&TasksToExecute[0], resultFile, nil)
		}

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
