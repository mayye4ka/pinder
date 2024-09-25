package service

import (
	"context"
	"fmt"
	"pinder/models"
	"pinder/server"
)

func (s *Service) UpdProfile(ctx context.Context, req *server.UpdProfileRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	mapped := mapProfile(req.NewProfile)
	mapped.UserID = userId
	if mapped.Gender != models.GenderFemale && mapped.Gender != models.GenderMale {
		return fmt.Errorf("bad profile")
	}
	if mapped.LocationName == "" || mapped.Name == "" {
		return fmt.Errorf("bad profile")
	}
	if mapped.Age == 0 || mapped.LocationLat == 0 || mapped.LocationLon == 0 {
		return fmt.Errorf("bad profile")
	}
	return s.repository.PutProfile(mapped)
}

func (s *Service) GetProfile(ctx context.Context, req *server.RequestWithToken) (*server.GetProfileResponse, error) {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	profile, err := s.repository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return &server.GetProfileResponse{}, nil
	}
	photoLinks, err := s.getUserPhotoLinks(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &server.GetProfileResponse{
		Profile: unmapProfile(*profile, photoLinks),
	}, nil
}

func (s *Service) UpdPreferences(ctx context.Context, req *server.UpdPreferencesRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	mapped := mapPreferences(req.NewPreferences)
	mapped.UserID = userId
	return s.repository.PutPreferences(mapped)
}

func (s *Service) GetPreferences(ctx context.Context, req *server.RequestWithToken) (*server.GetPreferencesResponse, error) {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	preferences, err := s.repository.GetPreferences(userId)
	if err != nil {
		return nil, err
	}
	if preferences == nil {
		return &server.GetPreferencesResponse{}, nil
	}
	return &server.GetPreferencesResponse{
		Preferences: unmapPreferences(*preferences),
	}, nil
}
