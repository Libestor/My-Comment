package cobalt_file

import (
	"fmt"
	"os"
)

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
func ReadFile(filename string) {
	file, _ := os.OpenFile(filename, os.O_RDONLY, 0400)
	var read_buffer = make([]byte, 10)
	var content_buffer = make([]byte, 0)
	fileinfo, _ := file.Stat()
	size := fileinfo.Size() //文件大小，单位是字节，int64
	var length int64 = 0    //标记已经读取了多少字节的内容
	for length < size {     //循环读取文件内容
		n, _ := file.Read(read_buffer)
		content_buffer = append(content_buffer, read_buffer[:n]...)
		length += int64(n)
	}
	fmt.Println(string(content_buffer))
}
