package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {
	searcher := Searcher{}
	searcher.maxChars = 300
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	PhraseIndexes []int

	maxChars int
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	sb := strings.Builder{}
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Can't read: %w", err)
		}
		sb.WriteString(strings.TrimSpace(string(line)))
		sb.WriteByte('\n')
	}
	s.CompleteWorks = sb.String()
	s.index()
	log.Printf("Loaded")
	return nil
}

func (s *Searcher) index() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	func() {
		defer wg.Done()
		s.SuffixArray = suffixarray.New([]byte(s.CompleteWorks))
		log.Printf("Suffix index done")
	}()
	func() {
		defer wg.Done()
		s.PhraseIndexes = indexPhrases(s.CompleteWorks)
		log.Printf("Sentence index done")
	}()
	wg.Wait()
}

func (s *Searcher) Search(query string) [][]string {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := [][]string{}
	for _, idx := range idxs {
		results = append(results, s.makeResult(idx, query))
	}
	return results
}

func (s *Searcher) makeResult(idx int, query string) []string {
	str := s.getTextAt(idx)
	res := highlightText(str, query)
	return strings.Split(res, "\n")
}

func (s *Searcher) getTextAt(idx int) string {
	pi := sort.Search(len(s.PhraseIndexes), func(i int) bool { return s.PhraseIndexes[i] > idx })
	from, to := s.selectPos(pi - 1)
	return s.CompleteWorks[from:to]
}

func (s *Searcher) selectPos(idx int) (int, int) {
	from, to := idx, idx+1
	l := len(s.PhraseIndexes)
	if to > (l - 1) {
		to = (l - 1)
	}
	changed := true
	for changed {
		changed = false
		if from > 0 && (s.PhraseIndexes[to]-s.PhraseIndexes[from-1]) <= s.maxChars {
			from--
			changed = true
		}
		if to < (l-1) && (s.PhraseIndexes[to+1]-s.PhraseIndexes[from]) <= s.maxChars {
			to++
			changed = true
		}
	}
	return s.PhraseIndexes[from], s.PhraseIndexes[to]
}

func highlightText(str, phr string) string {
	return strings.ReplaceAll(str, phr, "<b>"+phr+"</b>")
}

func indexPhrases(str string) []int {
	phraseIndicators := phraseEnds()
	res := make([]int, 0)
	l := len(str)
	res = append(res, 0)
	for i := 0; i < (l - 1); i++ {
		if phraseStart(str[i:i+2], phraseIndicators) {
			res = append(res, i+2)
			i = i + 1
		}
	}
	res = append(res, l)
	return res
}

func phraseEnds() map[string]bool {
	res := make(map[string]bool)
	for _, pe := range []string{"\n\n", ".\n", "?\n", "!\n", ". ", "? ", "! "} {
		res[pe] = true
	}
	return res
}

func phraseStart(str string, phrStrInd map[string]bool) bool {
	return phrStrInd[str]
}
