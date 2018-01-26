package service

import (
	"context"

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
	d, e := wordsoup.NewDictionary()
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
