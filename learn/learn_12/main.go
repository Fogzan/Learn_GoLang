package main

import "main/colorizer"

func main() {
	colorizer.PrintError("Ошибка!")
	colorizer.PrintInfo("Информация!")
	colorizer.PrintSuccess("Успешно!")
	colorizer.PrintWarning("Оповещение!")
}
