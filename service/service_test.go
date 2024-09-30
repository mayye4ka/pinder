package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testCtx     = context.Background()
	userId      = uint64(123)
	user2Id     = uint64(124)
	phoneNumber = "123456789"
	password    = "SuperMegaPassword123"
	passHash    = "a8630ec77eb54401e672f6fbb46c67304d02b9b747399b27f67e463b6878d7bb"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.vp68JTvxUceGTgtbZfHHato2w2Hjzuuv-Ne-Ts3kwcY"
	userName    = "name"

	photoKeys  = []string{"ph1", "ph2"}
	photoLinks = []string{"link1", "link2"}
)

type ServiceTestSuite struct {
	suite.Suite
	repoMock           *MockRepository
	fsMock             *MockFileStorage
	userInteractorMock *MockUserInteractor
	service            *Service
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repoMock = NewMockRepository(ctrl)
	s.fsMock = NewMockFileStorage(ctrl)
	s.userInteractorMock = NewMockUserInteractor(ctrl)
	s.service = New(s.repoMock, s.fsMock, s.userInteractorMock)
}
