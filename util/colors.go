package util

import (
	"fmt"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// An always incrementing index, so we have an idea what print was before what
var printIndex int = 0

func PrintWhite(str string) {
	PrintSpine(White + str + Reset)
}

func PrintGreen(str string) {
	PrintSpine(Green + str + Reset)
}

func PrintYellow(str string) {
	PrintSpine(Yellow + str + Reset)
}

func PrintBlue(str string) {
	PrintSpine(Blue + str + Reset)
}

func PrintPurple(str string) {
	PrintSpine(Purple + str + Reset)
}

func PrintRed(str string) {
	PrintSpine(Red + str + Reset)
}

func PrintSpine(str string) {
	currentTime := time.Now()

	// Get printIndex as a fixed width number
	printIndexStr := fmt.Sprintf("%04d", printIndex)

	fmt.Println(Gray + printIndexStr + " | " + currentTime.Format("15:04:05") + " | " + Reset + str)

	printIndex++
}
