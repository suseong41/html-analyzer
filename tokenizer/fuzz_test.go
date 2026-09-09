package tokenizer

import "testing"

// 토크나이저가 죽지 않고 끝내야함.
// 유효한 오프셋만 내놓는지 확인
func FuzzTokenizer(f *testing.F) {
	f.Add("<p>hello <b>world</b></p>")
	f.Add(`<img src=x onerror=alert(1)>`)
	f.Add("<script>if (a<b) x='</p>';</script>")
	f.Add("<script><!--<script></script>x//--></script>")
	f.Add("<!-- <script>alert(1)</script> -->")
	f.Add("<a href='&#106;avascript:alert&lpar;1&rpar;'>")
	f.Add("<!DOCTYPE html><title>a &lt; b</title>")
	f.Add("</></ b><?xml?><!bogus>")
	f.Add("<a href=")
	f.Add("<<<>>>")
	f.Add("안녕 <b>세계</b>")

	f.Fuzz(func(t *testing.T, s string) {
		z := New(s)
		for n := 0; ; n++ {
			tok := z.Next()
			if tok.Type == ErrToken {
				return
			}
			// 토큰 수는 입력 길이를 넘을 수 없다.
			if len(s)+16 < n {
				t.Fatalf("전진하지 않는 경로: %q", s)
			}
			// 오프셋은 입력 범위 안.
			if tok.Offset < 0 || len(s) < tok.Offset {
				t.Fatalf("토큰 오프셋 %d가 범위 밖 (len=%d)", tok.Offset, len(s))
			}
			// 위치 변환은 1-기반이며 죽지 않음.
			if line, col := z.Position(tok.Offset); line < 1 || col < 1 {
				t.Fatalf("Position(%d) = %d:%d", tok.Offset, line, col)
			}
			for _, a := range tok.Attrs {
				if a.Offset < 0 || len(s) < a.Offset {
					t.Fatalf("속성 오프셋 %d가 범위 밖 (len=%d)", a.Offset, len(s))
				}
			}
		}
	})
}

// Unescape는 상태 기계와 분리되어 따로 검사.
func FuzzUnescape(f *testing.F) {
	f.Add("&#106;avascript:")
	f.Add("&#xZZ;&#;&#x110000;&#0;")
	f.Add("?a=1&amp;b=2&copy=3")
	f.Add("&&&&&#####")

	f.Fuzz(func(t *testing.T, s string) {
		out := Unescape(s)
		if len(s) < len(out) { // 참조는 항상 줄이는 방향
			t.Fatalf("Unescape(%q)가 %d -> %d로 늘어남", s, len(s), len(out))
		}
	})
}
