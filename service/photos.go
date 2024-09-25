package service

import (
	"context"
	"pinder/server"
)

func (s *Service) getUserPhotoLinks(ctx context.Context, userId uint64) ([]string, error) {
	photos, err := s.repository.GetUserPhotos(userId)
	if err != nil {
		return nil, err
	}
	links := make([]string, len(photos))
	for i, photo := range photos {
		link, err := s.filestorage.MakeLink(ctx, photo)
		if err != nil {
			return nil, err
		}
		links[i] = link
	}
	return links, nil
}

func (s *Service) AddPhoto(ctx context.Context, token string, photo []byte) error {
	userId, err := verifyToken(token)
	if err != nil {
		return err
	}
	key, err := s.filestorage.SavePhoto(ctx, photo)
	if err != nil {
		return err
	}
	return s.repository.AddPhoto(userId, key)
}

func (s *Service) DeletePhoto(ctx context.Context, req *server.DelPhotoRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	err = s.filestorage.DelPhoto(ctx, req.PhotoKey)
	if err != nil {
		return err
	}
	return s.repository.DeleteUserPhoto(userId, req.PhotoKey)
}
