package cobalt_crypto

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

type Ecode interface {
}

func Encode(b []byte) (string, error) {

	return string(b), nil
}
func Decode(b []byte) (string, error) {
	a, err := GBKToUtf8(b)
	return string(a), err
}

func GBKToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
