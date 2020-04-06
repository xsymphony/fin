package main

import (
	"encoding/json"
	"net/http"

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
		c.String(http.StatusBadRequest, "param error")
		return
	}
	words, index := automaton.Find(req.Sentence)
	if len(index) == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data": map[string]interface{}{
				"words":    words,
				"replaced": "",
				"matched":  false,
			},
		})
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
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "ok",
		"data": map[string]interface{}{
			"replaced": string(replaced),
			"matched":  true,
		},
	})
}

func main() {
	automaton = ac.NewAutomaton()
	automaton.Add("暴力")
	automaton.Add("膜")
	automaton.Add("蛤")
	automaton.Build()

	r := fin.New()
	r.Apply(fin.HandleNotFound(func(c *fin.Context) {
		c.JSONAbort(http.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "not found",
		})
	}))
	{
		r.POST("/api/v1/replace", replaceWord)
	}
	r.Run(":8080")
}
