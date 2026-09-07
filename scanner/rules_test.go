package scanner

import (
	"testing"

	"github.com/suseong41/suseong-html-analyzer/tokenizer"
)

func TestRuleJavaScriptURL(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{`<a href="javascript:alert(1)">`, 1},
		{`<a href="JaVaScRiPt:alert(1)">`, 1},
		{`<a href="&#106;avascript:alert(1)">`, 1},
		{`<a href="java&#9;script:alert(1)">`, 1}, // 탭 우회
		{`<a href="javascript:;">`, 0},            // 오탐 제외
		{`<a href="javascript:void(0)">`, 0},      // 오탐 제외
		{`<a href="/normal.html">`, 0},
	}
	for _, c := range cases {
		tok := tokenizer.New(c.in).Next()
		if got := len(ruleJavaScriptURL(tok)); got != c.want {
			t.Errorf("%s → %d건, want %d건", c.in, got, c.want)
		}
	}
}
