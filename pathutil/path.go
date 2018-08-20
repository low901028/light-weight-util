package pathutil

import "path"

// CanonicalURLPath 返回规划的url路径, 规则:
// 1. 路径总是以 "/" 开始
// 2. 用单斜杠替换多个斜杠
// 3. 替换'.' '..'
// 4.保留尾部的斜杠
func CanonicalURLPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean 清除根以外的尾斜杠,
	// 必要时添加尾斜杠.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}
