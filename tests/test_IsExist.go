package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_IsExist() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_IsExist ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_IsExist (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_IsExist (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "config.ini"
	expect_output := true
	output := IsExist(input_1)

	fmt.Println("input_1：", input_1)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

	if output != expect_output {
		panic("结果不符合预期")
	}

}
