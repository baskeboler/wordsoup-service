package service

import (
	"context"
	"fmt"
	"os"

	"github.com/baskeboler/wordsoup"
)

// Service interface
type Service interface {
	//GetDictionary(c context.Context) (wordsoup.Dictionary, error)
	GeneratePuzzle(c context.Context, width, height, words int) (*wordsoup.WordSoup, error)
}

type serviceImpl struct {
	dict wordsoup.Dictionary
}

// NewService builds the service
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
	ws, e := wordsoup.GenerateRandomWordSoup(height, width, words, s.dict)
	if e != nil {
		return nil, e
	}
	return ws, nil
}
