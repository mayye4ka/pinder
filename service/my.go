package service

import "pinder/server"

func (s *Service) UpdProfile(req *server.UpdProfileRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	mapped := mapProfile(req.NewProfile)
	mapped.UserID = userId
	return s.repository.PutProfile(mapped)
}

func (s *Service) GetProfile(req *server.RequestWithToken) (*server.GetProfileResponse, error) {
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
	return &server.GetProfileResponse{
		Profile: unmapProfile(*profile),
	}, nil
}

func (s *Service) UpdPreferences(req *server.UpdPreferencesRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	mapped := mapPreferences(req.NewPreferences)
	mapped.UserID = userId
	return s.repository.PutPreferences(mapped)
}

func (s *Service) GetPreferences(req *server.RequestWithToken) (*server.GetPreferencesResponse, error) {
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
