package tools

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5(source string) string {
	h := md5.New()
	io.WriteString(h, source)
	return fmt.Sprintf("%x", h.Sum(nil))
}
