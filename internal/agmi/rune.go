package agmi

var (
	isVSpace    = isRune('\n') // We currently don't allow for '\r'
	isHSpace    = isRune(' ', '\t')
	notIsHSpace = notIsRune(' ', '\t')
)

func isRune(rs ...rune) func(rune) bool {
	return func(r1 rune) bool {
		for _, r2 := range rs {
			if r1 == r2 {
				return true
			}
		}
		return false
	}
}

func notIsRune(rs ...rune) func(rune) bool {
	f := isRune(rs...)
	return func(r1 rune) bool {
		return !f(r1)
	}
}
