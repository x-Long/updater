package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Read_config() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Read_config ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Read_config (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Read_config (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "config.ini"
	expect_output := map[string]string{
		"beta_hosts":      "test_1.com, 123.57.105.16",
		"check_url":       "http://123.57.105.167/check_version.json",
		"current_version": "1.4.1.3",
		"install_root":    "./localfile/install_root",
	}
	output := Read_config(input_1)

	fmt.Println("input_1：", input_1)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

}
