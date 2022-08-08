package cobalt_crypto

import (
	cobalt_file "My-Comment/cobalt.file"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

var DEBUG bool = false
var c1 int32
var c2 int32
var c3 int32

func getCnum() {
	var m = int(time.Now().Month())
	var d = time.Now().Day()
	var con1 int32 = 1
	var con2 int32 = 1
	var con3 int32 = 1
	m1 := (d + m) * int(math.Abs(float64(d-m)))
	m2 := (d << 1) * (m << 1) * 10000
	m3 := d * m * 1000000
	for m1/int(math.Pow(10, float64(con1))) > 0 {
		con1++
	}
	for m2/int(math.Pow(10, float64(con2))) > 0 {
		con2++
	}
	for m3/int(math.Pow(10, float64(con3))) > 0 {
		con3++
	}
	c1 = con1*10 + con1
	c2 = con2*10 + con2
	c3 = con3*10 + con3
}
func EncodeString(s string) (string, error) {
	getCnum()
	s = base64.StdEncoding.EncodeToString([]byte(s)) // base64加密
	r := []rune(s)
	newB := runeEncode(r)
	return string(newB), nil
}
func EncodeByte(b []byte) ([]byte, error) {
	getCnum()
	dis := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	if len(b) <= 2 {
		return []byte(""), nil
	}
	base64.StdEncoding.Encode(dis, b) //base64加密
	r := []rune(string(dis))
	newB := runeEncode(r)
	return []byte(string(newB)), nil
	//return b, nil
}
func DecodeToByte(b []byte) ([]byte, error) {
	getCnum()
	r := []rune(string(b))
	newB := runeDecode(r)
	newByte := []byte(string(newB))
	var err1 error = nil
	//var b []byte = '15'

	newByte = bytes.Trim(newByte, "\r")
	newByte = bytes.Trim(newByte, "\n")
	newByte = bytes.Trim(newByte, " ")
	i := 1
	for i <= 32 {
		newByte = bytes.Trim(newByte, string(30))
		i++
	}

	if bytes.Index(newByte, []byte("!@#$^&*()_+")) != -1 {
		num := bytes.Index(newByte, []byte("!@#$^&*()_+"))
		B1 := newByte[:num]
		B2 := newByte[num+11:]
		if DEBUG {
			fmt.Printf("分离前:")
			fmt.Println(newByte)
		}
		B1, err2 := base64.StdEncoding.DecodeString(string(B1)) //base64解密
		cobalt_file.PutErr(err2, "")
		if DEBUG {
			fmt.Printf("解密后B1：%s\n", B1)
			fmt.Printf("此时B2：%s\n", B2)
		}
		B2, err3 := base64.StdEncoding.DecodeString(string(B2)) //base64解密

		cobalt_file.PutErr(err3, "")
		if DEBUG {
			fmt.Printf("解密后B2：%s\n", B2)
		}
		B3 := bytes.Join([][]byte{B1, B2}, []byte("")) // 连接起来
		return B3, err3
	}
	if DEBUG {
		fmt.Printf("解密后未base:%s\n", newByte)
	}
	if DEBUG {
		fmt.Printf(" 二进制文件：")
		fmt.Println(newByte)
	}
	strs := string(newByte)

	newByte1, err1 := base64.StdEncoding.DecodeString(strs) //base64解密

	if DEBUG {
		fmt.Printf("内部解密后数据：%s\n", newByte1)
	}

	return newByte1, err1

}
func DecodeToString(b []byte) (string, error) {
	c, err := DecodeToByte(b)

	return string(c), err
}

func GBKToUtf8(s []byte) ([]byte, error) {

	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return d, e
	}

	return d, nil
}
func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return d, e
	}

	return d, nil
}
func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func runeDecode(old []rune) []rune {

	for i := 0; i < len(old); i++ {
		old[i] = coreDeCode(old[i])
	}

	return old
}
func runeEncode(old []rune) []rune {

	for i := 0; i < len(old); i++ {
		old[i] = coreEncode(old[i])
	}

	return old
}
func coreDeCode(rune2 rune) rune {
	rune2 = rune2 ^ c3
	rune2 = rune2 ^ c2
	rune2 = rune2 ^ c1
	return rune2

}
func coreEncode(rune2 rune) rune {
	rune2 = rune2 ^ c1
	rune2 = rune2 ^ c2
	rune2 = rune2 ^ c3
	return rune2
}
func CryphoDebug(msg bool) {
	DEBUG = msg
}
