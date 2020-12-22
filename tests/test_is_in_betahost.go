package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Is_in_betahost() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Is_in_betahost ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Is_in_betahost (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Is_in_betahost (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "http://123.57.105.167/audit-2.4.6.3.zip"
	input_2 := "test_1.com, 192.168.1.1"
	expect_output := false
	output := Is_in_betahost(input_1, input_2)

	fmt.Println("input_1：", input_1)
	fmt.Println("input_1：", input_2)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

	if output != expect_output {
		panic("结果不符合预期")
	}

}
