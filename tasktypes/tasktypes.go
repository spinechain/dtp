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

	_ "embed"
)

//go:embed latent_diffusion.bat
var ld_bat string

//go:embed latent_diffusion.sh
var ld_sh string

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

func RunLatentDiffusion(dataFolder string, taskType string, prompt string) error {

	var shellScriptName string
	util.PrintYellow("Running latent diffusion: " + taskType)

	os := runtime.GOOS
	switch os {
	case "windows":
		if taskType == "ld" {
			// check if linux or windows
			shellScriptName = "latent_diffusion.bat"
		}
	case "darwin":
		util.PrintRed("LD is not supported on mac")
	case "linux":
		if taskType == "ld" {
			// check if linux or windows
			shellScriptName = "latent_diffusion.sh"
		}
	default:
		fmt.Printf("%s.\n", os)
	}

	if shellScriptName == "" {
		util.PrintRed("Task type not supported: " + taskType)
		return errors.New("not supported")
	}

	var cmd *exec.Cmd
	if filepath.Ext(shellScriptName) == ".bat" {
		// run the batch file
		cmd = exec.Command("cmd.exe", "/C", filepath.Join(dataFolder, "scripts", shellScriptName))
	} else {
		// run the shell script
		cmd = exec.Command("sh", filepath.Join(dataFolder, "scripts", shellScriptName))
	}

	result, err := cmd.Output()
	if err != nil {
		util.PrintRed("Error running latent diffusion: " + err.Error())
		return err
	}

	// result to string
	resultString := string(result)

	// search for text in result
	if strings.Contains(resultString, "not a valid Win32") {
		util.PrintRed("The latent diffusion script is not a valid Win32 application")
	}

	return nil
}
