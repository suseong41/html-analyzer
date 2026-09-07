package tokenizer

import "testing"

func TestSelfClosing(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"<br/>", true},
		{"<br />", true},
		{"<br>", false},
		{"<br / >", false},
		{"<a href=x/>", false},
		{`<a href="x"/>`, true},
		{"<div/>", true},
	}
	for _, c := range cases {
		tok := New(c.in).Next()
		if tok.SelfClosing != c.want {
			t.Errorf("%q: SelfClosing=%v want %v", c.in, tok.SelfClosing, c.want)
		}
	}
}

func TestPosition(t *testing.T) {
	cases := []struct {
		in           string
		wantL, wantC int
	}{
		{"<p>x</p>", 1, 1},
		{"\n\n<p>", 3, 1},
		{"<p>안녕</p>", 1, 1},
	}
	for _, c := range cases {
		z := New(c.in)
		tok := z.Next()
		if tok.Type == TextToken {
			tok = z.Next()
		}
		l, col := z.Position(tok.Offset)
		if l != c.wantL || col != c.wantC {
			t.Errorf("%q: got L%d:C%d want L%d:C%d", c.in, l, col, c.wantL, c.wantC)
		}
	}
}

// 한글이 섞인 줄에서 칸이 글자 기준인지 확인
func TestPositionRuneColumn(t *testing.T) {
	z := New("<p>안녕</p>")
	z.Next()        // <p>
	z.Next()        // 안녕
	tok := z.Next() // </p>
	if _, col := z.Position(tok.Offset); col != 6 {
		t.Errorf("col=%d, want 6 (바이트 기준이면 10이 나온다)", col)
	}
}
