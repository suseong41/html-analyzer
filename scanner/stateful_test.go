package scanner

import "testing"

func TestDomainOf(t *testing.T) {
	cases := []struct{ in, want string }{
		{"https://www.jnu.ac.kr/", "www.jnu.ac.kr"},
		{"http://a.com:8080/x?y", "a.com"},
		{"//cdn.example.com/x.js", "cdn.example.com"},
		{"https://evil.com@bank.com/login", "bank.com"},
		{"https://bank.com@evil.com/login", "evil.com"},
		{"/relative/path", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := domainOf(c.in); got != c.want {
			t.Errorf("domainOf(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormOrigin(t *testing.T) {
	const page = "https://bank.example.com/login"
	cases := []struct {
		name string
		html string
		want int
	}{
		{"외부도메인+비밀번호", `<form action="https://evil.com/x"><input type=password></form>`, 1},
		{"상대경로", `<form action="/login"><input type=password></form>`, 0},
		{"같은도메인", `<form action="https://bank.example.com/x"><input type=password></form>`, 0},
		{"비밀번호없음", `<form action="https://evil.com/x"><input type=text></form>`, 0},
		{"form안닫힘", `<form action="https://evil.com/x"><input type=password>`, 1},
		{"대문자타입", `<form action="https://evil.com/x"><input TYPE=PASSWORD></form>`, 1},
		{"유저정보속임", `<form action="https://bank.example.com@evil.com/x"><input type=password></form>`, 1},
		{"폼밖의비밀번호", `<input type=password><form action="https://evil.com/x"></form>`, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var n int
			for _, f := range ScanURL(c.html, page).Findings {
				if f.Code == "cross-origin-password-form" {
					n++
				}
			}
			if n != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, n, c.want)
			}
		})
	}
}

func TestZeroWidth(t *testing.T) {
	cases := []struct {
		html string
		want int
	}{
		{"<p>정상 텍스트</p>", 0},
		{"<p>보이지\u200b않음</p>", 1},
		{"<a href=\"a\u200bb.com\">x</a>", 1},
		{"<p>\uFEFF</p>", 1},
	}
	for _, c := range cases {
		var n int
		for _, f := range Scan(c.html).Findings {
			if f.Code == "zero-width" {
				n++
			}
		}
		if n != c.want {
			t.Errorf("%q → %d건, want %d건", c.html, n, c.want)
		}

	}
}
