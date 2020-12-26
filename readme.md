# 一、自升级程序配置说明


## 1、准备

> linux家目录：  cd ~  
> windows 家目录一般为：cd C:/Users/用户名  

1. 依赖安装

```bash
cd home_path/go/src
mkdir -r golang.org/x
cd golang.org/x
git clone https://github.com/golang/text.git   
```
2. 项目拉取

```bash
cd home_path/go/src
git clone https://github.com/x-Long/updater
```

## 2、编译

```bash
cd home_path/go/src/updater
go build -ldflags="-w -s" updater.go	# 编译生成updater.exe
.\upx\upx.exe -9 updater.exe	# 对updater.exe 进行压缩
```

## 3、运行
 

- 若代码有改动，必须先进行上一步编译，然后使用  ```path/updater.exe [命令行参数]``` 的形式运行或测试
- **运行时注意 ```--install_root``` 为必填命令行参数**

```bash
./updater.exe --install_root C:\Users\long\Desktop\localfile\install_root
```
## 4、运行参数说明

```bash
--install_root (必填,安装目录)

--log_path (选填,日志路径)
--config (选填,配置文件目录)
--check_url (选填,版本json文件)

--current_version (选填，优先级最高)
--version_path (选填，优先级中等)

安装目录下version文件优先级最低
```
## 5、测试

```bash
cd home_path/go/src/updater       # 在updater目录下测试
go run .\test_updater.go
```

## 6、退出代码

```bash
exit(0)  # 升级成功
exit(1)  # 不用升级
exit(2)  # 升级失败
```

# 二、自升级程序需求说明

## 1、背景说明

假设你用 pyqt 制作编写了一款工具软件，使用 `pyinstaller` 集成依赖，基于 `NSIS` 生成安装包给客户端使用。每当有 bug 解决或者新增特性时，为了方便用户使用到最新版本的软件。需要一个程序，定期检测是否有新版本存在，并将其下载到本地。

## 2、开发要求

### 1、功能要求

1. 开发语言只能使用 Golang 
1. 使用 github 平台托管代码，命名为 `updater` 每天都要有 commit 记录
1. 每个模块需要有测试函数，具体要求参考此项目 [proxychains-ng](https://github.com/rofl0r/proxychains-ng/tree/master/tests)
1. 以 YYYY-mm-dd HH:MM:DD log-level 的格式，记录日志到文件中
1. 不要被网络异常，文件权限错误等情况搞得 crash，谨慎处理异常

### 2、功能要求

1. 支持配置文件和命令行的解析
2. 提供构建脚本，一键获取源码，安装依赖并自动构建可执行程序
3. 提供测试脚本，一条命令跑完所有单元测试

### 3、系统支持
-  win7/win10
-  ubuntu 18.04 

## 3、配置文件格式

```sh
check_url = http://example.com:8090/client/check_version

# 如果升级服务器位于以下主机名单中，直接升级，不计算升级概率
beta_hosts = test_1.com, 192.168.114.1
```

**技术注解**: 该配置文件与本程序(updater)平级放置名为 `config.ini` 程序启动时默认加载该配置文件。


## 4、命令行说明

```shell
# 如果没有提供此参数，输出到 stdout
--log-path /path/to/log

--current_version 2.4.6.2
--version_path /path/to/version # 该文件只有版本号独占一行，再无其他任何信息

--check_url 192.168.114.1/client/check_version
--instatlled_root /path/to/program/installed
# 若同时提供命令行参数与配置文件，以命令行参数为准，命令行中未提供的参数，从配置文件中读取

--config /path/to/config/file

--check_url http://example.com/client/check_version
```

## 5、升级流程

1. 客户端发送 HTTP GET 请求 `/client/check_version` 得到版本列表，每个版本字段如下
	```json
	{
		"os": "windows",
		"version": "2.4.6.3",
		"update_percent": 30,
		"md5": "111111111111jfldakjfskadjflsjfl",
		"url": "https://example/update/windows/audit-1.2.3.zip",
	}
	```
	**技术注解**: `update_percent` 表示有几成用户可以更新。先让一小部分人升级，如果没有人反馈问题，慢慢增大 update_percent 的数字，直到 100

	该请求应填写如下内容到 HTTP header

	'User-Agent': 'audit-client-{current_version} OS:{os_version} COMPUTER_NAME: {computer_name}'
	
2. 比较本地和远程版本，如果发现新版本，则进行如下计算
	```python
	# 以下是伪码,判断当前客户端是否被“选中”升级
	# 0 <= updater_percent <= 100
	def is_chosen_update(beta_hosts: List(str), update_percent: int) -> bool:
		# 从 check_url 字符串中解析出 check_host
	    if check_host in beta_hosts:
	        return True
		random_seed = get_home_path() # 以用户家目录为种子
	    seed_ascii_sum = sum([ord(i) for i in random_seed])
	    lucky_number = seed_ascii_sum % 100
	    can_update = lucky_number <= update_percent
	    return can_update
	```
	
3. 下载更新包
	1. 下载更新包到临时临时文件 `updater_package`
	2. 计算升级包哈希，如果 `md5` 与服务器返回的不一致，放弃本次升级 exit(2)
	3. 记录 `installed_root` 里已有的文件为 `exists_files`
	4. 解压升级包所有文件到一个临时目录 `tmp_extract_dir`，解压完毕后将`tmp_extract_dir` 里的所有内容移动到`installed_root`目录
	5. 删除升级包 `updater_package`,删除`exists_files`所有文件
	6. 成功升级 exit(0)

	