package tokenizer

import "testing"

/* go 에서 %q는 따옴표 감싼거, %v는 기본 형식
* want: 테스트 대상 함수가 반환해야 하는 올바른 기준값(기대값)
* got: 함수를 실제로 실행했을 때 나온 값(실제값)
 */
func TestUnescape(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain text", "plain text"},
		{"&#106;avascript:alert(1)", "javascript:alert(1)"},
		{"&#x6A;avascript:", "javascript:"},
		{"&#X6a;avascript:", "javascript:"},
		{"&#0000106;", "j"},
		{"&#106", "j"},
		{"java&#9;script:", "java\tscript:"},
		{"&Tab;javascript:", "\tjavascript:"},
		{"javascript&colon;alert&lpar;1&rpar;", "javascript:alert(1)"},
		{"a &amp; b", "a & b"},
		{"&lt;script&gt;", "<script>"},
		{"&notarealentity;", "&notarealentity;"},
		{"just & alone", "just & alone"},
		{"a&b=1", "a&b=1"},
		{"&#;", "&#;"},
		{"&#x;", "&#x;"},
		{"&#xZZ;", "&#xZZ;"},
		{"&#0;", "\uFFFD"},
		{"&#xD800;", "\uFFFD"},
		{"&#x110000;", "\uFFFD"},
		{"&#38;#106;", "&#106;"}, // 이중 인코딩: 한 번만
		{"?a=1&amp;b=2", "?a=1&b=2"},
	}
	for _, c := range cases {
		if got := Unescape(c.in); got != c.want {
			t.Errorf("Unescape(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRawVsRcdata(t *testing.T) {
	cases := []struct {
		in          string
		wantDecoded string
		wantRaw     bool
	}{
		{"<script>a &lt; b</script>", "a &lt; b", true},
		{"<style>a &lt; b</style>", "a &lt; b", true},
		{"<title>a &lt; b</title>", "a < b", false},
		{"<textarea>a &lt; b</textarea>", "a < b", false},
		{"<p>a &lt; b</p>", "a < b", false},
	}
	for _, c := range cases {
		z := New(c.in)
		z.Next()
		tok := z.Next()
		if tok.Raw != c.wantRaw || tok.DecodedData() != c.wantDecoded {
			t.Errorf("%q: Raw=%v Decoded=%q, want  Raw=%v Decoded=%q",
				c.in, tok.Raw, tok.DecodedData(), c.wantRaw, c.wantDecoded)
		}
	}
}

func TestDecodedAttr(t *testing.T) {
	tok := New(`<a href="&#106;avascript:alert&lpar;1&rpar;">`).Next()
	a := tok.Attrs[0]
	if a.Value != "&#106;avascript:alert&lpar;1&rpar;" {
		t.Errorf("원본이 보존되지 않음: %q", a.Value)
	}
	if got := a.DecodedValue(); got != "javascript:alert(1)" {
		t.Errorf("DecodedValue() = %q", got)
	}
}
