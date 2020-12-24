package main

import "os"
import "path"
import "path/filepath"
import . "updater/tools"
import "strings"

func main() {
	dir_now, _ := os.Getwd()
	dir_now = strings.Replace(dir_now, "\\", "/", -1)
	go_path := Getparentdirectory(dir_now)
	go_path = Getparentdirectory(go_path)
	os.Setenv("GOPATH", go_path)

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
	latest_version := Get_latest_update(Get_check_version(Check_url))
	is_update := Is_need_update(latest_version)

	if !is_update {
		InfoLogger.Printf("不用更新，准备退出...")
		os.Exit(1)

	} else {

		fileName := Download_update_file(path.Join(Install_root, ""), latest_version.Url)
		remote_file_hash := Get_hash(fileName)
		InfoLogger.Printf("远程文件md5：", remote_file_hash)
		if latest_version.Md5 == remote_file_hash {

			InfoLogger.Printf("md5校验成功")
			Zip_depress(fileName)
			InfoLogger.Printf("自动更新完成")
			os.Exit(0)

		} else {
			ErrorLogger.Printf("md5校验失败")
			os.Exit(2)
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
