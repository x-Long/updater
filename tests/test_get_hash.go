package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Get_hash() {

	Count = Count + 1

	fmt.Println("\n", "------------------RUN test_Get_hash ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Get_hash (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Get_hash (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "config.ini"
	output := Get_hash(input_1)
	fmt.Println("output：")
	fmt.Println(output)

}
