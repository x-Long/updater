package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Get_latest_update() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Get_latest_update ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Get_latest_update (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Get_latest_update (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := Get_check_version("http://123.57.105.167/check_version.json")
	output := Get_latest_update(input_1)

	// fmt.Println("input_1：",input_1)
	fmt.Println("output：", output)

}
