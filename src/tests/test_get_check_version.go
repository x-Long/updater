package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Get_check_version() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Get_check_version ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Get_check_version (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Get_check_version (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "http://123.57.105.167/check_version.json"
	output := Get_check_version(input_1)
	fmt.Println("output：")
	fmt.Println(output)

}
