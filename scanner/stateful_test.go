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
	const code = "cross-origin-password-form"
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
		{"중첩form", `<form action="https://evil.com/x"><form action="/safe"><input type=password></form></form>`, 1},
		{"div속깊이", `<form action="https://evil.com/x"><div><p><input type=password>`, 1},
		{"표속", `<form action="https://evil.com/x"><table><tr><td><input type=password>`, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, page, code); got != c.want {
				t.Errorf("%s → %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}

func TestZeroWidth(t *testing.T) {
	const code = "zero-width"
	cases := []struct {
		name string
		html string
		want int
	}{
		{"정상", "<p>정상 텍스트</p>", 0},
		{"텍스트", "<p>보이지\u200b않음</p>", 1},
		{"속성값", "<a href=\"a\u200bb.com\">x</a>", 1},
		{"BOM", "<p>\uFEFF</p>", 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := countCode(c.html, "", code); got != c.want {
				t.Errorf("%q → %d건, want %d건", c.html, got, c.want)
			}
		})
	}
}
