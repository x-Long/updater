package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_VersionGreaterThanT() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_VersionGreaterThanT ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_VersionGreaterThanT (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_VersionGreaterThanT (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "5.2.3.1"
	input_2 := "3.5.6.8"
	expect_output := true
	output := VersionGreaterThanT(input_1, input_2)

	fmt.Println("input_1：", input_1)
	fmt.Println("input_2：", input_2)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

	if output != expect_output {
		panic("结果不符合预期")
	}

}
