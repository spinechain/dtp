package tasktypes

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"spinedtp/util"
	"strings"

	_ "embed"
)

type TaskTypeSetting struct {
	LD_Script_Path string
}

var TaskSettings TaskTypeSetting

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

func RunLatentDiffusion(dataFolder string, shellScriptName string, prompt string) error {
	// Print
	util.PrintYellow("Running latent diffusion: " + shellScriptName)

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
