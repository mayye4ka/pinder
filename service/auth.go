package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"pinder/server"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("kefjhjdfh")

func getPassHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createToken(userId uint64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
		},
	)
	token, err := t.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func verifyToken(token string) (uint64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return 0, nil
	}
	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claimsMap, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("can't read claims")
	}
	id, ok := claimsMap["user_id"].(float64)
	if !ok {
		return 0, errors.New("can't read user_id")
	}
	return uint64(id), nil
}

func (s *Service) RegisterUser(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	passHash := getPassHash(req.Password)
	user, err := s.repository.CreateUser(req.PhoneNumber, passHash)
	if err != nil {
		return nil, err
	}
	token, err := createToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &server.RegisterResponse{Token: token}, nil
}

func (s *Service) LoginUser(ctx context.Context, req *server.LoginRequest) (*server.LoginResponse, error) {
	passHash := getPassHash(req.Password)
	user, err := s.repository.GetUserByCreds(req.PhoneNumber, passHash)
	if err != nil {
		return nil, err
	}
	token, err := createToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &server.LoginResponse{Token: token}, nil
}
