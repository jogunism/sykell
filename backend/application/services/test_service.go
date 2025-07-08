package services

import (
	"backend/application/queries"
)

type TestService struct{}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) GetTestMessage(query queries.GetTestMessageQuery) string {
	return "hello, world!"
}
