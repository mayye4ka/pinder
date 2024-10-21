package service

import (
	"context"
	"testing"

	"github.com/mayye4ka/pinder/internal/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	userId   = uint64(123)
	user1Ctx = context.WithValue(context.Background(), userIdContextKey, userId)
	user2Id  = uint64(124)

	userName   = "John"
	photo1     = "ph1"
	photo1Link = "link1"
	photo2     = "ph2"
	photo2Link = "link2"
	photos     = []models.PhotoShowcase{
		{
			Key:  photo1,
			Link: photo1Link,
		},
		{
			Key:  photo2,
			Link: photo2Link,
		},
	}
)

type ServiceTestSuite struct {
	suite.Suite
	repoMock         *MockRepository
	fsMock           *MockFileStorage
	userNotifierMock *MockUserNotifier
	sttMock          *MockStt
	service          *Service
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repoMock = NewMockRepository(ctrl)
	s.fsMock = NewMockFileStorage(ctrl)
	s.userNotifierMock = NewMockUserNotifier(ctrl)
	s.sttMock = NewMockStt(ctrl)
	s.service = New(s.repoMock, s.fsMock, s.userNotifierMock, s.sttMock)
}
