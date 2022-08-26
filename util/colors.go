package util

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func PrintGreen(str string) {
	fmt.Println(Green + str + Reset)
}

func PrintYellow(str string) {
	fmt.Println(Yellow + str + Reset)
}

func PrintPurple(str string) {
	fmt.Println(Purple + str + Reset)
}

func PrintRed(str string) {
	fmt.Println(Red + str + Reset)
}
