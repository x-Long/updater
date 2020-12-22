package tests

import "fmt"

var (
	Count int
	Pass  int
	Fail  int
)

func Test_updater() {

	test_Read_config()
	test_IsExist()
	test_Download_update_file()
	test_Get_http_header()
	test_Get_check_version()
	test_Get_latest_update()
	test_Is_in_betahost()
	test_Get_hash()
	test_VersionGreaterThanT()

	fmt.Println("\n", "-------统计----------- ")
	fmt.Println("Count", Count)
	fmt.Println("Pass", Pass)
	fmt.Println("Fail", Fail)

}
