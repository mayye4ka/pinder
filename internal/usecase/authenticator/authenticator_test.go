package authenticator

import (
	"context"
	"testing"

	"github.com/mayye4ka/pinder/internal/models"
	"github.com/rs/zerolog"
	gomock "go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testCtx     = context.Background()
	phoneNumber = "12345678"
	password    = "SuperMegaPassword123"
	passHash    = "a8630ec77eb54401e672f6fbb46c67304d02b9b747399b27f67e463b6878d7bb"
	userId      = uint64(123)
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.vp68JTvxUceGTgtbZfHHato2w2Hjzuuv-Ne-Ts3kwcY"
)

func TestGetPassHash(t *testing.T) {
	hash := getPassHash(password)
	assert.Equal(t, passHash, hash)
}

type AuthenticatorTestSuite struct {
	suite.Suite
	repoMock      *MockRepository
	authenticator *Authenticator
}

func (s *AuthenticatorTestSuite) TestRegisterUser() {
	s.repoMock.EXPECT().CreateUser(phoneNumber, passHash).Return(models.User{ID: userId}, nil)

	gotToken, err := s.authenticator.Register(testCtx, phoneNumber, password)

	s.Nil(err)
	s.Equal(token, gotToken)
}

func (s *AuthenticatorTestSuite) TestLoginUser() {
	s.repoMock.EXPECT().GetUserByCreds(phoneNumber, passHash).Return(models.User{ID: userId}, nil)

	gotToken, err := s.authenticator.Login(testCtx, phoneNumber, password)

	s.Nil(err)
	s.Equal(token, gotToken)
}

func (s *AuthenticatorTestSuite) TestUnpackToken() {
	id, err := s.authenticator.UnpackToken(testCtx, token)

	s.Nil(err)
	s.Equal(userId, id)
}

func (s *AuthenticatorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repoMock = NewMockRepository(ctrl)
	l := zerolog.Nop()
	s.authenticator = New(s.repoMock, &l)
}

func TestAuthenticator(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}
