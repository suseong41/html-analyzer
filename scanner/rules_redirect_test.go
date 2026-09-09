package scanner

import "testing"

func TestMetaRefreshScheme(t *testing.T) {
	const code = "meta-refresh-scheme"
	cases := []struct {
		name string
		html string
		want int
	}{
		// 음성 — 정상 페이지에서 나오면 안 되는 것들
		{"정상리다이렉트", `<meta http-equiv="refresh" content="0;url=/login">`, 0},
		{"외부정상", `<meta http-equiv="refresh" content="5; url=https://example.com/">`, 0},
		{"refresh아님", `<meta http-equiv="content-type" content="0;url=data:text/html,x">`, 0},
		{"content없음", `<meta http-equiv="refresh">`, 0},
		{"url없음", `<meta http-equiv="refresh" content="5">`, 0},
		// 양성
		{"dataURI", `<meta http-equiv="refresh" content="0;url=data:text/html;base64,PHNjcmlwdD4=">`, 1},
		{"javascript", `<meta http-equiv="refresh" content="0;url=javascript:alert(1)">`, 1},
		{"대문자", `<META HTTP-EQUIV="REFRESH" CONTENT="0;URL=DATA:text/html,x">`, 1},
		{"작은따옴표", `<meta http-equiv="refresh" content="0; url='data:text/html,x'">`, 1},
		{"공백많음", `<meta http-equiv="refresh" content="0 ;  url  =  data:text/html,x">`, 1},
		{"문자참조우회", `<meta http-equiv="refresh" content="0;url=&#100;ata:text/html,x">`, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, "https://a.com/", code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}

func TestBaseHrefExternal(t *testing.T) {
	const page = "https://a.com/x"
	const code = "base-href-external"
	cases := []struct {
		name, html, url string
		want            int
	}{
		{"base없음", `<p>x</p>`, page, 0},
		{"상대경로", `<base href="/app/">`, page, 0},
		{"같은도메인", `<base href="https://a.com/app/">`, page, 0},
		{"href없음", `<base target="_blank">`, page, 0},
		{"URL모름", `<base href="https://evil.com/">`, "", 0},
		{"외부도메인", `<base href="https://evil.com/">`, page, 1},
		{"스킴리스외부", `<base href="//evil.com/">`, page, 1},
		{"유저정보속임", `<base href="https://a.com@evil.com/">`, page, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, c.url, code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}

func TestIframeSandboxEscape(t *testing.T) {
	const code = "iframe-sandbox-escape"
	cases := []struct {
		name, html string
		want       int
	}{
		{"sandbox없음", `<iframe src="https://x.com/"></iframe>`, 0},
		{"빈sandbox", `<iframe src="https://x.com/" sandbox></iframe>`, 0},
		{"스크립트만", `<iframe src="https://x.com/" sandbox="allow-scripts"></iframe>`, 0},
		{"동일출처만", `<iframe src="https://x.com/" sandbox="allow-same-origin"></iframe>`, 0},
		{"둘다", `<iframe src="https://x.com/" sandbox="allow-scripts allow-same-origin"></iframe>`, 1},
		{"순서반대", `<iframe src="https://x.com/" sandbox="allow-same-origin allow-forms allow-scripts"></iframe>`, 1},
		{"대문자", `<iframe SANDBOX="ALLOW-SCRIPTS ALLOW-SAME-ORIGIN"></iframe>`, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, "https://a.com/", code); got != c.want {
				t.Errorf("%s -> %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}
