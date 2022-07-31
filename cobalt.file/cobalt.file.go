package cobalt_file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var Userpath = "UserOpersion"
var MaxMsgLen = 100

func OpenFIle(filename string) ([]byte, int64, error) {
	file, err1 := os.OpenFile(filename, os.O_RDONLY, 0400)
	if err1 != nil {
		return nil, 0, err1
	}
	var read_buffer = make([]byte, 100)
	var content_buffer = make([]byte, 0)
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}
	size := fileinfo.Size() //文件大小，单位是字节，int64
	var length int64 = 0    //标记已经读取了多少字节的内容
	for length < size {     //循环读取文件内容
		n, _ := file.Read(read_buffer)
		content_buffer = append(content_buffer, read_buffer[:n]...)
		length += int64(n)
	}
	return content_buffer, size, nil
}
func PutErr(err error, str string) {
	if err != nil {
		fmt.Printf(str)
		log.Println(err)
	}
}
func PathExists(paths string) bool {
	_, err := os.Stat(paths)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func MemUser(msg []byte) {
	msg1 := msg[:]
	if len(msg1) > MaxMsgLen {
		msg1 = msg1[:MaxMsgLen]
	}
	times := []byte(time.Now().Format("2006-01-02 15:04:05"))
	times = append(times, byte('\t'))
	msg1 = append(times, msg1...)

	name := time.Now().Format("01_02") + ".txt"
	files, err1 := os.OpenFile(Userpath+"/"+name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664) //二次打开使用
	PutErr(err1, "用户行为记录文件打开失败\n")
	write := bufio.NewWriter(files)
	_, err2 := write.Write(msg1)
	PutErr(err2, "用户行为写入失败\n")
	err3 := write.Flush()
	PutErr(err3, "用户操作文件到缓冲失败失败\n")
	err4 := files.Close()
	PutErr(err4, "文件关闭失败\n")
}
