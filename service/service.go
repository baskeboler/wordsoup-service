package service

import (
	"context"
	"fmt"
	"os"

	"github.com/baskeboler/wordsoup"
)

// type Dictionary wordsoup.Dictionary
// type WordSoup wordsoup.WordSoup
type Service interface {
	//GetDictionary(c context.Context) (wordsoup.Dictionary, error)
	GeneratePuzzle(c context.Context, width, height, words int) (*wordsoup.WordSoup, error)
}

type serviceImpl struct {
	dict wordsoup.Dictionary
}

func NewService() (Service, error) {
	dictPath := os.Getenv("WORDSOUP_DICT")
	fmt.Printf("Loading dictionary from %s\n", dictPath)
	d, e := wordsoup.NewDictionaryFromTextFile(dictPath)
	if e != nil {
		return nil, e
	}
	return &serviceImpl{dict: d}, nil
}

func (s *serviceImpl) GeneratePuzzle(c context.Context, width, height, words int) (*wordsoup.WordSoup, error) {
	// var ws interface{}
	ws, e := wordsoup.GenerateRandomWordSoup(height, width, words, s.dict)
	if e != nil {
		return nil, e
	}
	return ws, nil
}
