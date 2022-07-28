package cobalt_crypto

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

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
	r := []rune(s)
	newB := runeEncode(r)
	return string(newB), nil
}
func EncodeByte(b []byte) ([]byte, error) {
	getCnum()
	r := []rune(string(b))
	newB := runeEncode(r)
	return []byte(string(newB)), nil
	//return b, nil
}
func DecodeToByte(b []byte) ([]byte, error) {
	//Utf, _ := GBKToUtf8(b)
	getCnum()
	b, _ = GBKToUtf8(b)
	r := []rune(string(b))
	newB := runeDecode(r)
	newByte := []byte(string(newB))
	var err1 error = nil
	//newByte, err1 = GBKToUtf8([]byte(newByte))
	//newByte, err1 = zhToUnicode(newByte)
	//newByte, err1 = GBKToUtf8(newByte)
	//newByte, _ = UTF8ToGBK(newByte)
	return newByte, err1

}
func DecodeToString(b []byte) (string, error) {
	//b, _ = GBKToUtf8(b)
	c, err := DecodeToByte(b)
	if err != nil {
		return "", err
	}
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
func qer() {
	a := "中文字符"
	b := []rune(a)
	fmt.Println(b)
	var c []rune
	for _, j := range b {
		j += 2
		c = append(c, j)
		fmt.Println(j)
	}
}
func runeDecode(old []rune) []rune {

	//fmt.Println(b)
	//var c []rune
	for i := 0; i < len(old); i++ {
		old[i] = coreDeCode(old[i])
	}
	//for _, j := range old  {
	//	c = append(c, coreCode(j))
	//	//fmt.Println(j)
	//}
	return old
}
func runeEncode(old []rune) []rune {

	//fmt.Println(b)
	//var c []rune
	for i := 0; i < len(old); i++ {
		old[i] = coreEncode(old[i])
	}
	//for _, j := range old  {
	//	c = append(c, coreCode(j))
	//	//fmt.Println(j)
	//}
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
