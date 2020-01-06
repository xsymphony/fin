package main

import (
	"encoding/json"

	"github.com/xsymphony/ac"
	"github.com/xsymphony/fin"
)

var automaton *ac.Automaton

type replaceSensitiveRequest struct {
	Sentence string `json:"sentence"`
	Symbol   string `json:"symbol"`
}

func replaceWord(c *fin.Context) {
	var req replaceSensitiveRequest
	if err := json.Unmarshal(c.Request.Body(), &req); err != nil {
		c.SetStatusCode(499)
		c.WriteString("param error")
		return
	}
	c.Response.Header.Set("Content-type", "application/json")
	words, index := automaton.Find(req.Sentence)
	if len(index) == 0 {
		resp, _ := json.Marshal(map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data": map[string]interface{}{
				"words":    words,
				"replaced": "",
				"matched":  false,
			},
		})
		c.Write(resp)
		return
	}
	symbol := []rune(req.Symbol)[0]
	runes := []rune(req.Sentence)
	replaced := make([]rune, len(runes))
	var cursor int
	for i := 0; i < len(runes); i++ {
		start, end := index[cursor], index[cursor+1]
		if i >= start && i <= end {
			replaced[i] = symbol
		} else {
			replaced[i] = runes[i]
		}
		if i == end && cursor+1 < len(index)-1 {
			cursor += 2
		}
	}
	resp, _ := json.Marshal(map[string]interface{}{
		"code":    0,
		"message": "ok",
		"data": map[string]interface{}{
			"replaced": string(replaced),
			"matched":  true,
		},
	})
	c.Write(resp)
}

func main() {
	automaton = ac.NewAutomaton()
	automaton.Add("暴力")
	automaton.Add("膜")
	automaton.Add("蛤")
	automaton.Build()

	r := fin.New()
	{
		r.POST("/api/v1/replace", replaceWord)
	}
	r.Run(":8080")
}
