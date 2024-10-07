package authenticator

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
	"github.com/rs/zerolog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var jwtSecretKey = []byte("kefjhjdfh")

type Authenticator struct {
	repo   Repository
	logger *zerolog.Logger
}

type Repository interface {
	CreateUser(phoneNumber, passHash string) (models.User, error)
	GetUserByCreds(phoneNumber, passHash string) (models.User, error)
}

func New(repo Repository, logger *zerolog.Logger) *Authenticator {
	return &Authenticator{
		repo:   repo,
		logger: logger,
	}
}

func (a *Authenticator) Register(ctx context.Context, phone, password string) (string, error) {
	passHash := getPassHash(password)
	user, err := a.repo.CreateUser(phone, passHash)
	if err != nil {
		return "", errors.Wrap(err, "can't register new user")
	}
	token, err := a.createToken(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "can't generate token for registered user")
	}
	return token, nil
}

func (a *Authenticator) Login(ctx context.Context, phone, password string) (string, error) {
	passHash := getPassHash(password)
	user, err := a.repo.GetUserByCreds(phone, passHash)
	if err != nil {
		return "", errors.Wrap(err, "can't get user by creds")
	}
	token, err := a.createToken(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "can't generate token for logged in user")
	}
	return token, nil
}

func (a *Authenticator) UnpackToken(ctx context.Context, token string) (uint64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return 0, &errs.CodableError{
			Code:    errs.CodePermissionDenied,
			Message: "invalid token",
		}
	}
	if !t.Valid {
		return 0, &errs.CodableError{
			Code:    errs.CodePermissionDenied,
			Message: "invalid token",
		}
	}

	claimsMap, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, &errs.CodableError{
			Code:    errs.CodePermissionDenied,
			Message: "can't read claims",
		}
	}
	id, ok := claimsMap["user_id"].(float64)
	if !ok {
		return 0, &errs.CodableError{
			Code:    errs.CodePermissionDenied,
			Message: "can't read user_id",
		}
	}
	return uint64(id), nil
}

func getPassHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (a *Authenticator) createToken(userId uint64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
		},
	)
	token, err := t.SignedString(jwtSecretKey)
	if err != nil {
		a.logger.Err(err).Msg("can't create token")
		return "", &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create token",
		}
	}
	return string(token), nil
}
