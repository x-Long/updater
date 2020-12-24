package main

import "os"
import . "updater/tools"
import . "updater/tests"


func main() {

	Command_line()
	Test_updater()
	os.Exit(0)
}