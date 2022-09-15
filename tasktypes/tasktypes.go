package tasktypes

import (
	"log"
	"os/exec"
	"spinedtp/util"
	"strings"
)

type TaskTypeSetting struct {
	LD_Script_Path string
}

var TaskSettings TaskTypeSetting

func RunLatentDiffusion(shellScriptLocation string, prompt string) error {
	// Print
	util.PrintYellow("Running latent diffusion: " + shellScriptLocation)

	cmd := exec.Command(shellScriptLocation)

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

	if result != nil {
		log.Fatal(result)
	}
	log.Println(cmd.Run())

	return nil
}
