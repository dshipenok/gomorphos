package str

import "strings"

type Word []rune

func (w Word) Len() int {
	return len(w)
}

func (w Word) String() string {
	return string(w)
}

func (w Word) Lower() Word {
	return Word(
		strings.ToLower(
			string(w),
		),
	)
}

func (w Word) Upper() Word {
	return Word(
		strings.ToUpper(
			string(w),
		),
	)
}

func (w Word) LastChars(count int) string {
	return string(w.LastCharsWord(count))
}

func (w Word) LastCharsWord(count int) Word {
	len := w.Len()
	if count >= len {
		return w
	}
	return w[len-count:]
}

func (w Word) EndsWith(lastCount int, strs ...string) bool {
	end := w.LastChars(lastCount)
	for _, s := range strs {
		if end == s {
			return true
		}
	}
	return false
}

func (w Word) OneOf(strs ...string) bool {
	str := w.String()
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func (w Word) Chars(from, till int) string {
	return string(w.SliceWord(from, till))
}

func (w Word) SliceWord(from, till int) Word {
	len := w.Len()
	if from < 0 {
		from = len + from
	}
	if till < 0 {
		till = len + till
	}
	if till < from {
		till = from
	}
	if till > len {
		till = len
	}

	return w[from:till]
}

func (w Word) SubWord(from int) Word {
	len := w.Len()
	if from < 0 {
		from = len + from
	}
	if from >= len {
		return Word{}
	}

	return w[from:]
}

func (w Word) LastIndex(char string) int {
	for i := len(w) - 1; i >= 0; i-- {
		if string(w[i]) == char {
			return i
		}
	}
	return -1
}
