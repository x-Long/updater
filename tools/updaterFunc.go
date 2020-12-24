package tools

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	flag1      = log.Ldate | log.Ltime | log.Lshortfile
	preDebug   = "[DEBUG]"
	preInfo    = "[INFO]"
	preWarning = "[WARNING]"
	preError   = "[ERROR]"
)

var (
	logFile       io.Writer
	debugLogger   *log.Logger
	InfoLogger    *log.Logger
	warningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func Log_init() {
	var err error
	logFile, err = os.OpenFile(Log_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println("create log file err %+v", err)
		logFile = os.Stdout
	}
	debugLogger = log.New(logFile, preDebug, flag1)
	InfoLogger = log.New(logFile, preInfo, flag1)
	warningLogger = log.New(logFile, preWarning, flag1)
	ErrorLogger = log.New(logFile, preError, flag1)
}

func Debugf(format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func Warningf(format string, v ...interface{}) {
	warningLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// ---------------------------------------

var (
	Check_url       string
	Current_version string
	Install_root    string
	Log_path        string
	Is_test         bool
	Config          string
	Version_path    string
)

type Version_list struct {
	O_sys          string `json:"os"`
	Version        string `json:"version"`
	Update_percent int    `json:"update_percent"`
	Md5            string `json:"md5"`
	Url            string `json:"url"`
}

func Read_config(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Printf("read error:", "路径错误,请将配置文件config.ini放在程序同一目录下,或手动指定正确的目录")
		os.Exit(2)
		panic(err)
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read error:", "路径错误,请将配置文件config.ini放在程序同一目录下,或手动指定正确的目录")
			os.Exit(2)
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	Check_url = config["check_url"]
	// Install_root = config["install_root"]
	// Log_path = config["log_path"]
	// Current_version= config["current_version"]

	return config
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func Command_line() {

	if IsExist(path.Join(Install_root, "version")) {
		f, err := os.Open(path.Join(Install_root, "version"))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rd := bufio.NewReader(f)
		for {
			data, _, eof := rd.ReadLine()
			if eof == io.EOF {
				break
			}
			Current_version = string(data)

			break
		}
	}

	if Version_path != "" {
		if IsExist(Version_path) {
			f1, err1 := os.Open(Version_path)
			if err1 != nil {
				panic(err1)
			}
			defer f1.Close()

			rd1 := bufio.NewReader(f1)
			for {
				data1, _, eof1 := rd1.ReadLine()
				if eof1 == io.EOF {
					break
				}
				Current_version = string(data1)
				break
			}
		}
	}
	flag.StringVar(&Check_url, "check_url", Check_url, " ")
	flag.StringVar(&Install_root, "install_root", Install_root, " ")
	flag.StringVar(&Log_path, "log_path", Log_path, " ")
	flag.StringVar(&Current_version, "current_version", Current_version, " ")
	flag.StringVar(&Version_path, "version_path", "", " ")

	// flag.BoolVar(&Is_test, "test", false, " ")
	flag.StringVar(&Config, "config", Config, " ")
	flag.Usage = usage
	flag.Parse()

}

func usage() {
	updater_path, _ := filepath.Abs(os.Args[0])
	fmt.Println("Usage of " + updater_path + ":")
	fmt.Println("  " + "--help")
	fmt.Println("		" + "(帮助信息)")

	fmt.Println("  " + "--install_root")
	fmt.Println("		" + "(待升级程序的安装目录,若未指定,默认为updater.exe所在目录)")
	fmt.Println("  " + "--log_path")
	fmt.Println("		" + "(日志路径)")
	fmt.Println("  " + "--config")
	fmt.Println("		" + "(配置文件目录,若未指定,默认读取updater.exe同级目录下的config.ini文件)")
	fmt.Println("  " + "--check_url")
	fmt.Println("		" + "(获取版本json文件的url,若未指定,默认从config.ini中读取check_url字段)")
	fmt.Println("  " + "--current_version")
	fmt.Println("		" + "(当前版本,版本信息的优先级：current_version > version_path > 安装目录下的version文件)")
	fmt.Println("  " + "--version_path")
	fmt.Println("		" + "(记录当前版本的文件,版本信息的优先级：current_version > version_path > 安装目录下的version文件))")

}

func Download_update_file(Install_root, url string) string {
	res, err := http.Get(url)
	if err != nil {
		ErrorLogger.Printf("获取版本列表失败,请检查网络或Check_url")
		os.Exit(2)
		panic(err)
	}

	_, fileName := filepath.Split(url)
	InfoLogger.Printf("正在下载服务器数据...")
	os.RemoveAll(path.Join(filepath.Dir(path.Join(Install_root, "")), "zip_tmp"))
	os.MkdirAll(path.Join(filepath.Dir(path.Join(Install_root, "")), "zip_tmp"), os.ModePerm)

	fileName = path.Join(path.Join(filepath.Dir(path.Join(Install_root, "")), "zip_tmp"), fileName)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	f.Close()
	InfoLogger.Printf("服务器数据下载成功")
	return fileName
}

func Get_http_header() string {
	InfoLogger.Printf("生成http头部...")
	hostname, _ := os.Hostname()
	InfoLogger.Printf("http头部生成成功")
	return "audit-client-" + Current_version + " OS:" + runtime.GOOS + " COMPUTER_NAME:" + hostname

}

func Get_check_version(url string) string {
	InfoLogger.Printf("请求版本列表...")
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)

	reqest.Header.Add("User-Agent", Get_http_header())

	if err != nil {
		ErrorLogger.Printf("版本列表请求失败,请检查网络")
		os.Exit(2)
		panic(err)
	}
	res, err1 := client.Do(reqest)
	InfoLogger.Printf("%v", res.StatusCode)
	if err1 != nil {
		ErrorLogger.Printf("版本列表请求失败,请检查网络")
		os.Exit(2)
		panic(err)
	}

	defer res.Body.Close()

	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	InfoLogger.Printf("版本列表请求成功")
	return string(robots)
}

func Get_latest_update(remote_version_json string) Version_list {

	InfoLogger.Printf("开始获取最新版本...")
	var version_list []Version_list
	var jsonBlob = []byte(remote_version_json)
	err1 := json.Unmarshal(jsonBlob, &version_list)
	if err1 != nil {
		ErrorLogger.Printf("error:", "获取版本列表失败,请检查网络或Check_url")
		os.Exit(2)
		panic(err1)
	}
	InfoLogger.Printf("%s", version_list)

	latest_version := Current_version

	InfoLogger.Printf("查询最新版本...")
	for _, value := range version_list {
		if runtime.GOOS == value.O_sys {
			flag := VersionGreaterThanT(latest_version, value.Version)
			if flag {
				latest_version = latest_version
			} else {
				latest_version = value.Version
			}
		}
	}
	InfoLogger.Printf("最新版本", latest_version)

	index_out := -1

	for index, value := range version_list {
		if latest_version == value.Version {
			index_out = index
			break
		}
	}

	InfoLogger.Printf("最新版本信息", version_list[index_out])
	return version_list[index_out]
}

func Is_in_betahost(urls, data string) bool {
	InfoLogger.Printf("正在检查主机位置...")
	s := urls
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	// h := strings.Split(u.Host, ":")
	// log.Println(h[0])	//主机
	// log.Println(h[1])  //端口

	return strings.Contains(data, u.Host)

}

func Is_need_update(latest_version Version_list) bool {

	need := false

	if latest_version.Version == Current_version {
		InfoLogger.Printf("已经是最新版本,准备退出")
		return need
	}

	fd, err := os.Open(Config)

	defer fd.Close()

	if err != nil {
		ErrorLogger.Printf("read error:", "路径错误,请将配置文件config.ini放在程序同一目录下,或手动指定正确的目录")
		os.Exit(2)
		panic(err)

	}
	buff := bufio.NewReader(fd)

	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		if strings.Contains(string(data), "beta_hosts") {
			if Is_in_betahost(latest_version.Url, string(data)) {
				need = true
				break
			} else {
				InfoLogger.Printf("检查percent_rate...")
				// 检查家目录
				user, _ := user.Current()

				r := []uint8(user.HomeDir)
				var b int
				for _, value := range r {

					b += int(value)
				}
				lucky_number := b % 100
				can_update := lucky_number <= latest_version.Update_percent
				InfoLogger.Printf("更新比例：", latest_version.Update_percent, " 本地种子：", lucky_number)
				need = can_update
			}
		}
	}
	return need
}

func Get_hash(src string) string {

	InfoLogger.Printf("文件md5校验...")
	file, err := os.Open(src)
	defer file.Close()
	if err != nil {
		ErrorLogger.Printf("读取文件失败！")
	}
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)

	return hex.EncodeToString(sum)
}

func Zip_depress(fileName string) {

	InfoLogger.Printf("开始解压...")
	r1, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, k := range r1.Reader.File {

		if k.Flags == 0 {
			i := bytes.NewReader([]byte(k.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			k.Name = string(content)
		} else {
			// decodeName = f.Name
		}

		if k.FileInfo().IsDir() {
			err := os.MkdirAll(path.Join(path.Join(Install_root, ""), k.Name), os.ModePerm)

			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		r, err := k.Open()
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer r.Close()
		InfoLogger.Printf("正在提取: ", k.Name)

		fileName_zip := path.Join(path.Join(Install_root, ""), k.Name)

		NewFile, err := os.Create(fileName_zip)
		if err != nil {
			log.Fatal(err)
			continue
		}
		io.Copy(NewFile, r)
		NewFile.Close()
	}

	err = r1.Close()
	if err != nil {
		ErrorLogger.Printf("close file err=", err)
	}

	err1 := os.RemoveAll(path.Join(filepath.Dir(path.Join(Install_root, "")), "zip_tmp"))
	if err1 != nil {
		log.Fatal(err1)

	}
	InfoLogger.Printf("解压完成")
}

func VersionGreaterThanT(a, b string) bool {
	return GreaterThan(a, b)
}

func stripMetadata(v string) string {
	split := strings.Split(v, "+")
	if len(split) > 1 {
		return split[0]
	}
	return v
}

func periodDashSplit(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '.', '-':
			return true
		}
		return false
	})
}

func GreaterThan(a, b string) bool {

	numberRe := regexp.MustCompile("[0-9]+")
	wordRe := regexp.MustCompile("[a-z]+")
	a = stripMetadata(a)
	b = stripMetadata(b)

	a = strings.TrimLeft(a, "v")
	b = strings.TrimLeft(b, "v")

	aSplit := periodDashSplit(a)
	bSplit := periodDashSplit(b)

	if len(bSplit) > len(aSplit) {
		return !GreaterThan(b, a) && a != b
	}

	for i := 0; i < len(aSplit); i++ {
		if i == len(bSplit) {
			if _, err := strconv.Atoi(aSplit[i]); err == nil {
				return true
			}
			return false
		}
		aWord := wordRe.FindString(aSplit[i])
		bWord := wordRe.FindString(bSplit[i])
		if aWord != "" && bWord != "" {
			if strings.Compare(aWord, bWord) > 0 {
				return true
			}
			if strings.Compare(bWord, aWord) > 0 {
				return false
			}
		}
		aMatch := numberRe.FindString(aSplit[i])
		bMatch := numberRe.FindString(bSplit[i])
		if aMatch == "" || bMatch == "" {
			if strings.Compare(aSplit[i], bSplit[i]) > 0 {
				return true
			}
			if strings.Compare(bSplit[i], aSplit[i]) > 0 {
				return false
			}
		}
		aNum, _ := strconv.Atoi(aMatch)
		bNum, _ := strconv.Atoi(bMatch)
		if aNum > bNum {
			return true
		}
		if bNum > aNum {
			return false
		}
	}

	return false
}
