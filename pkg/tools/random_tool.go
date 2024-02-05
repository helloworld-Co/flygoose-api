package tools

import (
	"github.com/google/uuid"
	"math/big"
	"strings"
)
import crand "crypto/rand"

var characterArray = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func GenNumber(maxVal int64) int64 {
	n, _ := crand.Int(crand.Reader, big.NewInt(maxVal))
	return n.Int64()
}

// GenNumberCode 随机生成count个数字
func GenNumberCode(count int) string {
	nums := make([]string, count)
	for i := 0; i < count; i++ {
		n, _ := crand.Int(crand.Reader, big.NewInt(10))
		s := n.String()
		nums[i] = s
	}

	return strings.Join(nums, "")
}

func GenSectionId() int64 {
	num, _ := crand.Int(crand.Reader, big.NewInt(40000000))
	return num.Int64() + 542541
}

func GenSpecialId() int64 {
	num, _ := crand.Int(crand.Reader, big.NewInt(50000000))
	return num.Int64() + 7654321
}

func GenBlogId() int64 {
	num, _ := crand.Int(crand.Reader, big.NewInt(50000000))
	return num.Int64() + 12345678
}

// GenToken 随机生成token
func GenToken() string {
	return Md5(uuid.NewString())
}
