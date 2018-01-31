package service

import (
	"context"

	"github.com/baskeboler/auth"
	"github.com/baskeboler/wordsoup"
)

// Service interface
type Service interface {
	LoginService
	//GetDictionary(c context.Context) (wordsoup.Dictionary, error)
	GeneratePuzzle(c context.Context, width, height, words int) (*wordsoup.WordSoup, error)
}
type LoginService interface {
	Login(c context.Context, name string, password string) (*auth.LoginKey, error)
}
type serviceImpl struct {
	dict wordsoup.Dictionary
	auth auth.Manager
}

// NewService builds the service
func NewService() (Service, error) {
	d, e := wordsoup.NewDictionary()
	if e != nil {
		return nil, e
	}
	a := auth.New()

	a.CreateUser("user", "user")
	return &serviceImpl{dict: d, auth: a}, nil
}

func (s *serviceImpl) GeneratePuzzle(c context.Context, width, height, words int) (*wordsoup.WordSoup, error) {
	ws, e := wordsoup.GenerateRandomWordSoup(height, width, words, s.dict)
	if e != nil {
		return nil, e
	}
	return ws, nil
}
func (s *serviceImpl) Login(c context.Context, name string, password string) (*auth.LoginKey, error) {
	k, e := s.auth.GetKey(name, password)
	if e != nil {
		return nil, e
	}
	return k, nil
}
