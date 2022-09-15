package tasktypes

import (
	"log"
	"os"
	"os/exec"
	"spinedtp/util"
)

type TaskTypeSetting struct {
	LD_Script_Path string
}

var TaskSettings TaskTypeSetting

func RunLatentDiffusion(shellScriptLocation string, prompt string) {
	// Print
	util.PrintYellow("Running latent diffusion: " + shellScriptLocation)

	cmd := exec.Command(shellScriptLocation, prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Run())
}
