package cobalt_crypto

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math"
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
	c1 = con1*1000 + con1*100 + con1*10 + con1
	c2 = con2*1000 + con2*100 + con2*10 + con2
	c3 = con3*1000 + con3*100 + con3*10 + con3
}
func EncodeString(s string) (string, error) {
	getCnum()
	r := []rune(string(s))
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
func DecodeToString(b []byte) (string, error) {
	//Utf, _ := GBKToUtf8(b)
	getCnum()
	r := []rune(string(b))
	newB := runeDecode(r)
	newByte := string(newB)
	newByte, err1 := GBKToUtf8([]byte(newByte))
	return newByte, err1

}
func DecodeToByte(b []byte) ([]byte, error) {
	c, err := DecodeToString(b)
	if err != nil {
		return []byte(""), err
	}
	return []byte(c), err
}

func GBKToUtf8(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return string(d), e
	}

	return string(d), nil
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
