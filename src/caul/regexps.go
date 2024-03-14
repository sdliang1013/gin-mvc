package caul

import "regexp"

var (
	cameCaseCompiler = regexp.MustCompile("[A-Z][a-z]*")
)

// SplitCameCase 拆分驼峰字符串
func SplitCameCase(str string) []string {
	return cameCaseCompiler.FindAllString(str, -1)
}
