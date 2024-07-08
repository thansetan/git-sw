package gitconfig

type char interface {
	rune | byte
}

func isAlpha[T char](ch T) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNum[T char](ch T) bool {
	return ch >= '0' && ch <= '9'
}

func isAlnum[T char](ch T) bool {
	return isAlpha(ch) || isNum(ch)
}
