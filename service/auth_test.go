package service

import (
	"pinder/models"
	"pinder/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPassHash(t *testing.T) {
	hash := getPassHash(password)
	assert.Equal(t, passHash, hash)
}

func TestCreateToken(t *testing.T) {
	token, err := createToken(userId)
	assert.Nil(t, err)
	assert.Equal(t, token, token)
}

func TestVerifyToken(t *testing.T) {
	uid, err := verifyToken(token)
	assert.Nil(t, err)
	assert.Equal(t, userId, uid)
}

func (s *ServiceTestSuite) TestRegisterUser() {
	s.repoMock.EXPECT().CreateUser(phoneNumber, passHash).Return(&models.User{ID: userId}, nil)

	resp, err := s.service.RegisterUser(testCtx, &server.RegisterRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	})

	s.Nil(err)
	s.Equal(token, resp.Token)
}

func (s *ServiceTestSuite) TestLoginUser() {
	s.repoMock.EXPECT().GetUserByCreds(phoneNumber, passHash).Return(&models.User{ID: userId}, nil)

	resp, err := s.service.LoginUser(testCtx, &server.LoginRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	})

	s.Nil(err)
	s.Equal(token, resp.Token)
}
