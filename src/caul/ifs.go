package caul

func IfString(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func IfInt(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}

func IfInt32(cond bool, a, b int32) int32 {
	if cond {
		return a
	}
	return b
}

func IfInt64(cond bool, a, b int64) int64 {
	if cond {
		return a
	}
	return b
}

func IfError(cond bool, a, b error) error {
	if cond {
		return a
	}
	return b
}

func DefaultString(a, b string) string {
	return IfString(a != "", a, b)
}

func DefaultInt(a, b int) int {
	return IfInt(a != 0, a, b)
}

func DefaultInt32(a, b int32) int32 {
	return IfInt32(a != 0, a, b)
}

func DefaultInt64(a, b int64) int64 {
	return IfInt64(a != 0, a, b)
}

func DefaultError(a, b error) error {
	return IfError(a != nil, a, b)
}
