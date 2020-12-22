package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Get_http_header() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Get_http_header ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Get_http_header (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Get_http_header (", elapsed, ")")
		Pass = Pass + 1
	}()

	output := Get_http_header()
	fmt.Println("output：", output)

}
