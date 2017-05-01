package wordfilter

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type wordMap struct {
	lenMap   map[int]int
	lenSlice []int
	words    map[string]string
}

type wordTree struct {
	wordMaxLen int
	trees      map[string]*wordMap
}

func (t *wordTree) add(word string) {
	wSlice := strings.Split(word, "")
	length := len(wSlice)
	if length > t.wordMaxLen {
		t.wordMaxLen = length
	}
	if t.trees == nil {
		t.trees = make(map[string]*wordMap)
	}
	if t.trees[wSlice[0]] == nil {
		t.trees[wSlice[0]] = &wordMap{lenMap: make(map[int]int), words: make(map[string]string)}
	}
	tree := t.trees[wSlice[0]]
	if _, exist := tree.words[word]; !exist {
		tree.words[word] = strings.Repeat("*", length)
		if _, exist := tree.lenMap[length]; !exist {
			tree.lenMap[length] = length
			tree.lenSlice = append(tree.lenSlice, length)
			sort.Ints(tree.lenSlice)
		}
	}
}

type search struct {
	txt         *[]string
	txtLen      int
	tree        *wordTree
	replacement string
	matched     map[string]int
}

func (s *search) replace(start, end int) {
	for i := start; i < end; i++ {
		(*s.txt)[i] = s.replacement
	}
}

func (s *search) run() {
	s.matched = make(map[string]int)
	s.txtLen = len(*s.txt)
	for i := 0; i < s.txtLen; i++ {
		if tree, exist := s.tree.trees[(*s.txt)[i]]; exist {
			lenLen := len(tree.lenSlice)
			for j := 0; j < lenLen; j++ {
				end := i + tree.lenSlice[j]
				if end > s.txtLen {
					break
				}
				word := strings.Join((*s.txt)[i:end], "")
				if tree.words[word] != "" {
					if s.replacement != "" {
						s.replace(i, end)
					}
					sum, _ := s.matched[word]
					s.matched[word] = sum + 1
					i = end - 1
					break
				}
			}
		}
	}
}

func Init(filename string) *wordTree {
	wT := &wordTree{}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		wT.add(word)
	}
	return wT
}

func Search(tree *wordTree, txt *[]string, replacement string) *map[string]int {
	s := &search{}
	s.tree = tree
	s.txt = txt
	s.replacement = replacement
	s.run()
	return &s.matched
}

