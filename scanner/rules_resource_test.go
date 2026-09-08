package scanner

import (
	"strings"
	"testing"
)

func TestMixedContent(t *testing.T) {
	const https = "https://a.com/"
	const code = "mixed-content"
	cases := []struct {
		name, html, page string
		want             int
	}{
		{"링크는아님", `<a href="http://x.com/">z</a>`, https, 0},
		{"이미지", `<img src="http://x.com/a.png">`, https, 1},
		{"스타일시트", `<link rel=stylesheet href="http://x.com/a.css">`, https, 1},
		{"폼전송", `<form action="http://x.com/p"></form>`, https, 1},
		{"https리소스", `<img src="https://x.com/a.png">`, https, 0},
		{"상대경로", `<img src="/a.png">`, https, 0},
		{"URL모름", `<img src="http://x.com/a.png">`, "", 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, c.page, code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}

	t.Run("URL모름이면심각도낮춤", func(t *testing.T) {
		f, ok := findFirst(`<img src="http://x.com/a.png">`, "", code)
		if !ok {
			t.Fatal("발견되지 않음")
		}
		if f.Severity != Info {
			t.Errorf("심각도 = %v, want INFO", f.Severity)
		}
	})

	t.Run("https페이지면Medium", func(t *testing.T) {
		f, ok := findFirst(`<img src="http://x.com/a.png">`, https, code)
		if !ok {
			t.Fatal("발견되지 않음")
		}
		if f.Severity != Medium {
			t.Errorf("심각도 = %v, want MEDIUM", f.Severity)
		}
	})
}

func TestSubresourceIntegrity(t *testing.T) {
	const page = "https://a.com/"
	const code = "sri-missing"
	cases := []struct {
		name, html string
		want       int
	}{
		{"외부스크립트", `<script src="//cdn.x.com/a.js"></script>`, 1},
		{"integrity있음", `<script src="//cdn.x.com/a.js" integrity="sha384-aaa"></script>`, 0},
		{"상대경로", `<script src="/local.js"></script>`, 0},
		{"같은도메인", `<script src="https://a.com/x.js"></script>`, 0},
		{"인라인", `<script>var x = 1;</script>`, 0},
		{"스타일시트", `<link rel=stylesheet href="//cdn.x.com/a.css">`, 1},
		{"파비콘", `<link rel=icon href="//cdn.x.com/f.ico">`, 0},
		{"프리커넥트", `<link rel=preconnect href="//cdn.x.com/">`, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, page, code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}

func TestTargetBlank(t *testing.T) {
	const page = "https://a.com/"
	const code = "target-blank-no-rel"
	cases := []struct {
		name, html string
		want       int
	}{
		{"rel있음", strings.Repeat(`<a target=_blank rel=noopener href="/x">z</a>`, 3), 0},
		{"noreferrer도OK", `<a target=_blank rel="noreferrer nofollow" href="/x">z</a>`, 0},
		{"세곳이지만한건", strings.Repeat(`<a target=_blank href="/x">z</a>`, 3), 1},
		{"target없음", `<a href="/x">z</a>`, 0},
		{"대문자", `<a TARGET="_BLANK" href="/x">z</a>`, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, page, code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}

	t.Run("집계가세고있는가", func(t *testing.T) {
		html := strings.Repeat(`<a target=_blank href="/x">z</a>`, 3)
		f, ok := findFirst(html, page, code)
		if !ok {
			t.Fatal("발견되지 않음")
		}
		if !strings.Contains(f.Evidence, "3곳") {
			t.Errorf("Evidence = %q, 3곳이 세어지지 않음", f.Evidence)
		}
	})
}
