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
	DataFolder  string
	Prompt      string
	Complete    bool
	ResultFile  string
	ResultError error
}

// List with tasks to execute
var tasksToExecute []TaskToExecute

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

func AddToLatentDiffusionQueue(dataFolder string, taskType string, prompt string) error {
	util.PrintYellow("Adding to latent diffusion queue: " + taskType)

	// add to the queue
	tasksToExecute = append(tasksToExecute, TaskToExecute{TaskType: taskType, DataFolder: dataFolder, Prompt: prompt})

	// add to queue
	if !isRunning {
		go RunLatentDiffusion()
	} else {
		util.PrintYellow("Latent diffusion is already running")
	}

	return nil
}

func CompleteTask(task *TaskToExecute, resultFile string, resultError string) {
	task.Complete = true
	task.ResultFile = resultFile
	task.ResultError = errors.New(resultError)
}

// Get next open task
func GetNextTask() *TaskToExecute {
	for i, task := range tasksToExecute {
		if !task.Complete {
			return &tasksToExecute[i]
		}
	}

	return nil
}

func RunLatentDiffusion() error {

	isRunning = true
	defer func() { isRunning = false }()

	util.PrintYellow("Starting latent diffusion")

	// TODO: find the next open task

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
			CompleteTask(&tasksToExecute[0], "", "task type not supported")
			continue
		}

		var cmd *exec.Cmd
		if filepath.Ext(shellScriptName) == ".bat" {
			// run the batch file
			cmd = exec.Command("cmd.exe", "/C", filepath.Join(te.DataFolder, "scripts", shellScriptName), outputDir)
		} else {
			// run the shell script
			cmd = exec.Command("bash", filepath.Join(te.DataFolder, "scripts", shellScriptName), outputDir)
		}

		startTime := time.Now()

		result, err := cmd.Output()
		if err != nil {
			// Complete task
			CompleteTask(&tasksToExecute[0], "", err.Error())
			continue
		}

		// result to string
		resultString := string(result)

		// search for text in result
		if strings.Contains(resultString, "not a valid Win32") {
			util.PrintRed("The latent diffusion script is not a valid Win32 application")
		}

		if strings.Contains(resultString, "Enjoy.") {
			util.PrintBlue("Latent diffusion completed successfully")

			resultFile, err := FigureOutResultFile(outputDir, startTime)
			if err != nil {
				// Complete task
				CompleteTask(&tasksToExecute[0], "", err.Error())
				continue
			}

			// Complete task
			CompleteTask(&tasksToExecute[0], resultFile, "")
		}
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
