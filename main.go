package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/suseong41/suseong-html-analyzer/scanner"
)

// map: [키]값{}

func main() {

	// HTML READ
	data, err := os.ReadFile("testdata/jnu_main.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------- 처리 시작 --------")
	res := scanner.Scan(string(data))

	for _, f := range res.Findings {
		fmt.Printf("  %4d:%-4d %-4s [%s] %s\n", f.Line, f.Col, f.Severity, f.Code, f.Evidence)
	}
	fmt.Printf("\n탐지 %d건\n", len(res.Findings))
	fmt.Println("토큰 종류:", res.Tokens)
	names := make([]string, 0, len(res.Tags))
	for name := range res.Tags {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return res.Tags[names[j]] < res.Tags[names[i]]
	})

	for i, name := range names {
		if 15 <= i {
			break
		}
		fmt.Printf("  %-12s %d\n", name, res.Tags[name])
	}

	fmt.Println("-------- 처리 완료 --------")
}
