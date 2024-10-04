package authenticator

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/mayye4ka/pinder/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("kefjhjdfh")

type Authenticator struct {
	repo Repository
}

type Repository interface {
	CreateUser(phoneNumber, passHash string) (models.User, error)
	GetUserByCreds(phoneNumber, passHash string) (models.User, error)
}

func New(repo Repository) *Authenticator {
	return &Authenticator{
		repo: repo,
	}
}

func (a *Authenticator) Register(ctx context.Context, phone, password string) (string, error) {
	passHash := getPassHash(password)
	user, err := a.repo.CreateUser(phone, passHash)
	if err != nil {
		return "", err
	}
	token, err := createToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *Authenticator) Login(ctx context.Context, phone, password string) (string, error) {
	passHash := getPassHash(password)
	user, err := a.repo.GetUserByCreds(phone, passHash)
	if err != nil {
		return "", err
	}
	token, err := createToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *Authenticator) UnpackToken(ctx context.Context, token string) (uint64, error) {
	userId, err := unpackToken(token)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

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

func unpackToken(token string) (uint64, error) {
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
