package colorizer

import "github.com/fatih/color"

func PrintSuccess(msg string) {
	color.Green(msg)
}

func PrintError(msg string) {
	color.Red(msg)
}

func PrintWarning(msg string) {
	color.Yellow(msg)
}

func PrintInfo(msg string) {
	color.Blue(msg)
}
