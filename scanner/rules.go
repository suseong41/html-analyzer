package scanner

import (
	"strings"

	"github.com/suseong41/suseong-html-analyzer/tokenizer"
)

// normalizeURL() 브라우저가 실제로 보는 URL 값을 만든다.
// 디코딩 -> 소문자 -> 공백류 제거.
func normalizeURL(v string) string {
	v = tokenizer.Unescape(v)                 // &#106; → j
	v = strings.ToLower(strings.TrimSpace(v)) // JaVaScRiPt: → javascript:
	return strings.Map(func(r rune) rune {    // java\tscript: → javascript:
		if r == '\t' || r == '\n' || r == '\r' || r == '\f' {
			return -1
		}
		return r
	}, v)
}

// isDangerousJSURL() 실제로 코드가 실행되는 javascript: URL만 골라낸다.
func isDangerousJSURL(v string) bool {
	if !strings.HasPrefix(v, "javascript:") {
		return false
	}
	body := strings.Trim(strings.TrimPrefix(v, "javascript:"), "; ")
	return body != "" && body != "void(0)"
}

func ruleInlineHandler(tok tokenizer.Token) []Finding {
	var out []Finding
	for _, a := range tok.Attrs {
		if strings.HasPrefix(a.Name, "on") && 2 < len(a.Name) {
			out = append(out, Finding{
				Code:     "inline-handler",
				Title:    "인라인 이벤트 핸들러",
				Severity: Low,
				Offset:   a.Offset,
				Evidence: "<" + tok.Name + " " + a.Name + "=…>",
			})
		}
	}
	return out
}

func ruleJavaScriptURL(tok tokenizer.Token) []Finding {
	var out []Finding
	for _, a := range tok.Attrs {
		if a.Name != "href" && a.Name != "src" {
			continue
		}
		if isDangerousJSURL(normalizeURL(a.Value)) {
			out = append(out, Finding{
				Code:     "javascript-url",
				Title:    "javascript: URL",
				Severity: Medium,
				Offset:   a.Offset,
				Evidence: a.Name + "=" + a.Value,
			})
		}
	}
	return out
}
