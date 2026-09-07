package main

import (
	"fmt"
	"htmlscanner/tokenizer"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	test1 = "<p>hello <b>world</b></p>"
)

func main() {

	// HTML READ
	data, err := os.ReadFile("jnu_main.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------- 처리 시작 --------")
	z := tokenizer.New(string(data))
	// map: [키]값{}
	counts := map[tokenizer.TokenType]int{}
	tagCounts := map[string]int{}

	for {
		tok := z.Next()
		if tok.Type == tokenizer.ErrToken {
			break
		}
		counts[tok.Type]++
		if tok.Type == tokenizer.StartTagToken {
			for _, a := range tok.Attrs {
				if strings.HasPrefix(a.Name, "on") && 2 < len(a.Name) {
					l, c := z.Position(a.Offset)
					fmt.Printf("  %4d:%-3d [인라인 핸들러] <%s %s=…>\n", l, c, tok.Name, a.Name)
				}
				if (a.Name == "href" || a.Name == "src") &&
					strings.HasPrefix(strings.ToLower(strings.TrimSpace(a.Value)), "javascript:") {
					fmt.Printf("  [javascript: URL] <%s %s=%q>\n", tok.Name, a.Name, a.Value)
				}
			}
		}
	}
	fmt.Println("토큰 종류:", counts)

	names := make([]string, 0, len(tagCounts))
	for name := range tagCounts {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return tagCounts[names[j]] < tagCounts[names[i]] // 개수 내림차순
	})

	for i, name := range names {
		if 15 <= i {
			break
		}
		fmt.Printf("  %-12s %d\n", name, tagCounts[name])
	}

	fmt.Println("-------- 처리 완료 --------")
}
