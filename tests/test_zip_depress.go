package tests

import (
	"fmt"
	"time"
)

import . "updater/tools"

func test_Zip_depress() {

	Count = Count + 1
	fmt.Println("\n", "------------------RUN test_Zip_depress ")
	Log_init()
	start := time.Now()
	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Println("exception", exception)
			elapsed := time.Since(start)
			fmt.Println("------------------FAIL test_Zip_depress (", elapsed, ")")
			Fail = Fail + 1
			return

		}
		elapsed := time.Since(start)
		fmt.Println("------------------PASS test_Zip_depress (", elapsed, ")")
		Pass = Pass + 1
	}()

	input_1 := "./localfile/zip_tmp/audit-2.4.6.3.zip"
	Zip_depress(input_1)

}
