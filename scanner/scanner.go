package scanner

import "github.com/suseong41/suseong-html-analyzer/tokenizer"

type Severity int
type rule func(tok tokenizer.Token) []Finding

const (
	Info Severity = iota
	Low
	Medium
	High
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "INFO"
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	}
	return "?"
}

// Finding 정의하고 []Finding 슬라이스 사용.
type Finding struct {
	Code     string
	Title    string
	Severity Severity
	Offset   int    // Rule
	Line     int    // Scan()
	Col      int    // Scan()
	Evidence string // 원본
}

type Result struct {
	Findings []Finding
	Tokens   map[tokenizer.TokenType]int
	Tags     map[string]int
}

var rules = []rule{
	ruleInlineHandler,
	ruleJavaScriptURL,
}

func Scan(src string) Result {
	z := tokenizer.New(src)
	res := Result{
		Tokens: map[tokenizer.TokenType]int{},
		Tags:   map[string]int{},
	}
	for {
		tok := z.Next()
		if tok.Type == tokenizer.ErrToken {
			break
		}
		res.Tokens[tok.Type]++

		if tok.Type != tokenizer.StartTagToken {
			continue
		}
		res.Tags[tok.Name]++

		for _, r := range rules {
			res.Findings = append(res.Findings, r(tok)...) // []Finding 슬라이스를 하나씩 풀어서 전달.
		}
	}

	for i := range res.Findings {
		res.Findings[i].Line, res.Findings[i].Col = z.Position(res.Findings[i].Offset)
	}
	return res
}
