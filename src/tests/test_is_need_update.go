package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Is_need_update() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN Is_need_update ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL Is_need_update (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS Is_need_update (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := Get_latest_update(Get_check_version("http://123.57.105.167/check_version.json"))
	expect_output := true
	output := Is_need_update(input_1)

	// fmt.Println("input_1：",input_1)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

	if output != expect_output {
		panic("结果不符合预期")
	}

}
