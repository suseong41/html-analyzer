package tokenizer

import (
	"strings"
	"testing"
)

func TestScriptEscaped(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"평범한스크립트", `<script>var x = 1;</script>`,
			[]string{"START:script", "TEXT:var x = 1;", "END:script"}},

		{"주석만", `<script><!-- alert(1) --></script>`,
			[]string{"START:script", "TEXT:<!-- alert(1) -->", "END:script"}},

		{"escaped상태에서닫기", `<script><!-- a</script>`,
			[]string{"START:script", "TEXT:<!-- a", "END:script"}},

		{"이중이스케이프", `<script><!--<script></script>alert(1)//--></script>`,
			[]string{"START:script", `TEXT:<!--<script></script>alert(1)//-->`, "END:script"}},

		{"이중이스케이프안의태그", `<script><!--<script><img src=x onerror=alert(1)>//--></script>`,
			[]string{"START:script", `TEXT:<!--<script><img src=x onerror=alert(1)>//-->`, "END:script"}},

		{"주석닫고나서닫기", `<script><!--<script>a</script>--></script>`,
			[]string{"START:script", `TEXT:<!--<script>a</script>-->`, "END:script"}},

		{"문자열속주석", `<script>var s = "<!--";</script>`,
			[]string{"START:script", `TEXT:var s = "<!--";`, "END:script"}},

		{"style은영향없음", `<style><!--<script></script>--></style>`,
			[]string{"START:style", `TEXT:<!--<script></script>-->`, "END:style"}},
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
