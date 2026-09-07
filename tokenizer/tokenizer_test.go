package tokenizer

import (
	"strings"
	"testing"
)

func collect(input string) []string {
	z := New(input)
	var out []string
	for {
		tok := z.Next()
		if tok.Type == ErrToken {
			return out
		}
		s := tok.Data
		if tok.Type == StartTagToken || tok.Type == EndTagToken {
			s = tok.Name
			for _, a := range tok.Attrs {
				s += "|" + a.Name + "=" + a.Value
			}
		}
		out = append(out, tok.Type.String()+":"+s)
	}
}

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "기본",
			input: "<p>hello <b>world</b></p>",
			want:  []string{"START:p", "TEXT:hello ", "START:b", "TEXT:world", "END:b", "END:p"},
		},
		{
			name:  "한글",
			input: "안녕 <b>세계</b>",
			want:  []string{"TEXT:안녕 ", "START:b", "TEXT:세계", "END:b"},
		},
		{
			name:  "안닫힌태그",
			input: "<p>안 닫힘",
			want:  []string{"START:p", "TEXT:안 닫힘"},
		},
		{
			name:  "부등호",
			input: "<p>a < b</p>",
			want:  []string{"START:p", "TEXT:a < b", "END:p"},
		},
		{
			name:  "하트",
			input: "I <3 Go",
			want:  []string{"TEXT:I <3 Go"},
		},
		{
			name:  "끝의부등호",
			input: "abc<",
			want:  []string{"TEXT:abc<"},
		},
		{
			name:  "태그삼킴",
			input: "<p>a < <img src=x onerror=alert(1)></p>",
			want:  []string{"START:p", "TEXT:a < ", "START:img|src=x|onerror=alert(1)", "END:p"},
		},
		{
			name:  "대소문자",
			input: "<ScRiPt>x</ScRiPt>",
			want:  []string{"START:script", "TEXT:x", "END:script"},
		},
		{
			name:  "속성있음",
			input: `<img src=x onerror=alert(1)>`,
			want:  []string{"START:img|src=x|onerror=alert(1)"},
		},
		{
			name:  "슬래시로끝",
			input: "<br/>",
			want:  []string{"START:br"},
		},
		{
			name:  "닥타입",
			input: "<!DOCTYPE html>",
			want:  []string{"COMMENT:<!DOCTYPE html>"},
		},
		{
			name:  "스크립트내부",
			input: `<script>if (a<b) x = "</p>";</script>`,
			want:  []string{"START:script", `TEXT:if (a<b) x = "</p>";`, "END:script"},
		},
		{
			name:  "빈스크립트",
			input: "<script></script>",
			want:  []string{"START:script", "END:script"},
		},
		{
			name:  "닫기대소문자",
			input: "<script>x</SCRIPT>",
			want:  []string{"START:script", "TEXT:x", "END:script"},
		},
		{
			name:  "스타일",
			input: "<style>a > b { color: red }</style>",
			want:  []string{"START:style", "TEXT:a > b { color: red }", "END:style"},
		},
		{
			name:  "안닫힌스크립트",
			input: "<script>var x = 1;",
			want:  []string{"START:script", "TEXT:var x = 1;"},
		},
		{"따옴표3종", `<a href="a.html" title='hi' target=_blank>`,
			[]string{"START:a|href=a.html|title=hi|target=_blank"}},
		{"값없는속성", "<input disabled required>",
			[]string{"START:input|disabled=|required="}},
		{"속성대문자", "<IMG SRC=x ONERROR=alert(1)>",
			[]string{"START:img|src=x|onerror=alert(1)"}},
		{"따옴표안공백", `<a href="a b.html">`,
			[]string{"START:a|href=a b.html"}},
		{"등호주변공백", `<a href = "x">`,
			[]string{"START:a|href=x"}},
		{"탭개행구분", "<img\tsrc=x\nonerror=y>",
			[]string{"START:img|src=x|onerror=y"}},
		{"빈값", "<a href=>",
			[]string{"START:a|href="}},
		{"중복속성", `<a href="safe" href="javascript:x">`,
			[]string{"START:a|href=safe|href=javascript:x"}},
		{"닫는따옴표없음", `<a href="abc`,
			[]string{"START:a|href=abc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collect(tt.input)
			if strings.Join(got, "|") != strings.Join(tt.want, "|") {
				t.Errorf("\ninput: %q\ngot:  %q\nwant: %q", tt.input, got, tt.want)
			}
		})
	}
}
