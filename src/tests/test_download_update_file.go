package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Download_update_file() {
	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Download_update_file ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Download_update_file (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Download_update_file (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "./localfile/install_root"
	input_2 := "http://123.57.105.167/audit-2.4.6.3.zip"
	expect_output := "localfile/zip_tmp/audit-2.4.6.3.zip"
	output := Download_update_file(input_1, input_2)

	fmt.Println("input_1：", input_1)
	fmt.Println("expect_output：", expect_output)
	fmt.Println("output：", output)

	if output != expect_output {
		panic("结果不符合预期")
	}

}
