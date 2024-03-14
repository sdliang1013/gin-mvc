package caul

func ContainsString(ary []string, str string) bool {
	for _, s := range ary {
		if s == str {
			return true
		}
	}
	return false
}

func IndexString(ary []string, str string) int {
	for idx, s := range ary {
		if s == str {
			return idx
		}
	}
	return -1
}
