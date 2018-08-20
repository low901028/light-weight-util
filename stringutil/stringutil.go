package stringutil

import "math/rand"

const (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// UniqueStrings 返回随机生成的唯一字符串片段.
func UniqueStrings(maxlen uint, n int) []string {
	exist := make(map[string]bool)
	ss := make([]string, 0)

	for len(ss) < n {
		s := randomString(maxlen)
		if !exist[s] {
			exist[s] = true
			ss = append(ss, s)
		}
	}

	return ss
}

// RandomStrings 返回随机生成的字符串片段.
func RandomStrings(maxlen uint, n int) []string {
	ss := make([]string, 0)
	for i := 0; i < n; i++ {
		ss = append(ss, randomString(maxlen))
	}
	return ss
}

func randomString(l uint) string {
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}

