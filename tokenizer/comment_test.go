package tokenizer

import (
	"strings"
	"testing"
)

func dump(input string) []string {
	z := New(input)
	var out []string
	for {
		tok := z.Next()
		if tok.Type == ErrToken {
			return out
		}
		s := tok.Data
		if tok.Type == StartTagToken || tok.Type == EndTagToken || tok.Type == DoctypeToken {
			s = tok.Name
		}
		out = append(out, tok.Type.String()+":"+s)
	}
}

func TestComment(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"기본주석", "<!-- hello -->", []string{"COMMENT: hello "}},
		{"주석속태그", "<!-- <script>alert(1)</script> -->", []string{"COMMENT: <script>alert(1)</script> "}},
		{"빈주석", "<!-->", []string{"COMMENT:"}},
		{"빈주석2", "<!--->", []string{"COMMENT:"}},
		{"느낌표종료", "<!-- a --!>", []string{"COMMENT: a "}},
		{"안닫힌주석", "<!-- oops", []string{"COMMENT: oops"}},
		{"주석앞뒤텍스트", "a<!--b-->c", []string{"TEXT:a", "COMMENT:b", "TEXT:c"}},
		{"닥타입", "<!DOCTYPE html>", []string{"DOCTYPE:html"}},
		{"닥타입소문자", "<!doctype HTML>", []string{"DOCTYPE:html"}},
		{"가짜선언", "<!foo>", []string{"COMMENT:foo"}},
		{"xml선언", `<?xml version="1.0"?>`, []string{`COMMENT:?xml version="1.0"?`}},
		{"빈종료태그", "<p></></p>", []string{"START:p", "END:p"}},
		{"잘못된종료태그", "</ b>", []string{"COMMENT: b"}},
		{"종료태그EOF", "abc</", []string{"TEXT:abc", "TEXT:</"}},
		{"대시연속", "<!-- a --- b -->", []string{"COMMENT: a --- b "}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := dump(c.in)
			if strings.Join(got, "@") != strings.Join(c.want, "@") {
				t.Errorf("\ninput: %q\ngot:  %q\nwant: %q", c.in, got, c.want)
			}
		})
	}
}
