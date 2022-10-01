package util

import (
	"fmt"
	"log"
	"strings"
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

	toPrint := Gray + printIndexStr + " | " + currentTime.Format("15:04:05") + " | " + Reset + str
	fmt.Println(toPrint)

	toLog := printIndexStr + " | " + str

	toLog = strings.Replace(toLog, Reset, "", -1)
	toLog = strings.Replace(toLog, Yellow, "", -1)
	toLog = strings.Replace(toLog, Green, "", -1)
	toLog = strings.Replace(toLog, White, "", -1)
	toLog = strings.Replace(toLog, Red, "", -1)
	toLog = strings.Replace(toLog, Purple, "", -1)
	toLog = strings.Replace(toLog, Blue, "", -1)
	log.Println(toLog)
	printIndex++
}
