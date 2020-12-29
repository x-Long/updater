package main

import "os"
import "path"
import "path/filepath"
import . "tools"
import "strings"
import "time"
import "strconv"
import "math"

func main() {

	for index, value := range os.Args {
		if value == "--config" {
			Config = os.Args[index+1]
		}
		if value == "--version_path" {
			Version_path = os.Args[index+1]
		}
		if value == "--install_root" {
			Install_root = os.Args[index+1]
		}
		if value == "--check_url" {
			Check_url = os.Args[index+1]
		}
		if value == "--interval_check" {
			Interval_check = os.Args[index+1]
		}
	}
	if Install_root == "" {
		Install_root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	if Config == "" {
		Config, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		Config = Config + "/config.ini"
	}

	Read_config(Config)
	Command_line()
	Log_init()

	InfoLogger.Printf("Check_url=", string(Check_url), " Install_root=", string(path.Join(Install_root, "")))
	InfoLogger.Printf("Log_path=", string(Log_path), " Current_version=", string(Current_version))

	if !IsExist(Install_root) {
		ErrorLogger.Printf("查找安装目录失败")
		os.Exit(2)
	}

	hours, _ = strconv.ParseFloat(Interval_check, 64)
	if hours > 0 {
		for true {
			Interval_updater()
		}
	} else {
		Protect_run(Run_main_updater)
	}
}

var hours float64

func Interval_updater() {
	defer func() {
		err := recover()
		InfoLogger.Println(err)
		if err == "更新失败" {
			InfoLogger.Println("正在重试...")
			time.Sleep(time.Duration(10) * time.Second)
			return
		}
		time.Sleep(time.Duration(math.Floor(hours*3600)) * time.Second)
	}()
	Run_main_updater()
}

func Protect_run(entry func()) {
	defer func() {
		err := recover()
		switch err {
		case "成功更新":
			os.Exit(0)
		case "不用更新":
			os.Exit(1)
		case "更新失败":
			os.Exit(2)
		}
	}()
	entry()
}

func Run_main_updater() {

	latest_version := Get_latest_update(Get_check_version(Check_url))
	is_update := Is_need_update(latest_version)
	if !is_update {
		InfoLogger.Printf("不用更新，准备退出...")
		// os.Exit(1)
		panic("不用更新")

	} else {

		fileName := Download_update_file(path.Join(Install_root, ""), latest_version.Url)
		remote_file_hash := Get_hash(fileName)
		InfoLogger.Printf("远程文件md5：", remote_file_hash)
		if latest_version.Md5 == remote_file_hash {

			InfoLogger.Printf("md5校验成功")
			Zip_depress(fileName)

			InfoLogger.Printf("自动更新完成")
			Current_version = latest_version.Version
			// os.Exit(0)

			panic("成功更新")
		} else {
			ErrorLogger.Printf("md5校验失败")
			// os.Exit(2)
			panic("更新失败")
		}
	}
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Getparentdirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
