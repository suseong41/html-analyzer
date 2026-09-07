package tokenizer

import (
	"strings"
	"unicode/utf8"
)

/*
* rune 타입은 int32와 동일. 유니코드 코드 포인트
 */

/*
* 우회에 자주 쓰이는 것들
 */
var namedRefs = map[string]rune{
	"amp": '&', "lt": '<', "gt": '>', "quot": '"', "apos": '\'',
	"nbsp": '\u00A0', "Tab": '\t', "NewLine": '\n',
	"colon": ':', "semi": ';', "sol": '/', "bsol": '\\',
	"lpar": '(', "rpar": ')', "period": '.', "comma": ',',
	"excl": '!', "quest": '?', "equals": '=', "num": '#',
	"dollar": '$', "percnt": '%', "ast": '*', "plus": '+',
	"lowbar": '_', "grave": '`', "verbar": '|',
	"lbrace": '{', "rbrace": '}', "lsqb": '[', "rsqb": ']',
	"copy": '©', "reg": '®', "hellip": '…',
	"mdash": '—', "ndash": '–',
}

/*
* DecodedValue 블라우저가 실제로 보는 속성 값
* DecodedData 텍스트 토큰의 내용을 돌려줌.
* <script>/<style> 내부는 JS, CSS 소스이므로 해석하지 않음.
 */

func (a Attribute) DecodedValue() string { return Unescape(a.Value) }
func (tok Token) DecodedData() string {
	if tok.Raw {
		return tok.Data
	}
	return Unescape(tok.Data)
}

func Unescape(s string) string {
	if !strings.Contains(s, "&") {
		return s
	}
	var b strings.Builder // 내부 버퍼에 재할당 없이 이어붙임.
	b.Grow(len(s))        // 최종 크기 미리 잡아두기
	for i := 0; i < len(s); {
		if s[i] != '&' {
			b.WriteByte(s[i])
			i++
			continue
		}
		r, n := decodeRef(s[i:])
		if n == 0 {
			b.WriteByte('&')
			i++
			continue
		}
		b.WriteRune(r)
		i += n
	}
	return b.String()
}

func decodeRef(s string) (rune, int) {
	if len(s) < 2 || s[0] != '&' {
		return 0, 0
	}
	if s[1] == '#' {
		return decodeNumericRef(s)
	}
	j := 1
	for j < len(s) && isAlnum(s[j]) {
		j++
	}
	if j == 1 || len(s) <= j || s[j] != ';' {
		return 0, 0
	}
	if r, ok := namedRefs[s[1:j]]; ok {
		return r, j + 1
	}
	return 0, 0
}

func decodeNumericRef(s string) (rune, int) {
	i := 2
	base := 10
	if i < len(s) && (s[i] == 'x' || s[i] == 'X') {
		i++
		base = 16
	}
	start := i
	v := 0
	for i < len(s) && isDigitIn(s[i], base) {
		if v < 0x110000 {
			v = v*base + hexVal(s[i])
		}
		i++
	}
	if i == start {
		return 0, 0
	}
	if i < len(s) && s[i] == ';' {
		i++
	}
	return sanitizeCodePoint(v), i
}

func isAlnum(c byte) bool {
	return isASCIIAlpha(c) || ('0' <= c && c <= '9')
}

func hexVal(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10
	default:
		return int(c-'A') + 10
	}
}

func isDigitIn(c byte, base int) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if base == 16 {
		return ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
	}
	return false
}

func sanitizeCodePoint(v int) rune {
	if v == 0 || 0x10FFFF < v || (0x0800 <= v && v <= 0xDFFF) {
		return utf8.RuneError
	}
	return rune(v)
}
