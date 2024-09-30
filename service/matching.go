package service

import (
	"context"
	"pinder/models"
	"pinder/server"
)

func (s *Service) submitNextPartner(ctx context.Context, candidate *models.Profile) (*server.NextPartnerResponse, error) {
	links, err := s.getUserPhotoLinks(ctx, candidate.UserID)
	if err != nil {
		return nil, err
	}
	return &server.NextPartnerResponse{
		Partner: unmapProfile(*candidate, links),
	}, nil
}

func (s *Service) NextPartner(ctx context.Context, req *server.RequestWithToken) (*server.NextPartnerResponse, error) {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return nil, err
	}

	candidate, err := s.repository.GetHangingPartner(userId)
	if err != nil {
		return nil, err
	}
	if candidate != nil {
		return s.submitNextPartner(ctx, candidate)
	}

	candidate, err = s.repository.GetWhoLikedMe(userId)
	if err != nil {
		return nil, err
	}

	if candidate != nil {
		pa, err := s.repository.GetLatestPairAttempt(candidate.UserID, userId)
		if err != nil {
			return nil, err
		}
		err = s.repository.CreateEvent(pa.ID, models.PETypeSentToUser2)
		if err != nil {
			return nil, err
		}
		return s.submitNextPartner(ctx, candidate)
	}

	candidate, err = s.repository.ChooseCandidateAndCreatePairAttempt(userId)
	if err != nil {
		return nil, err
	}
	return s.submitNextPartner(ctx, candidate)
}

func (s *Service) Swipe(ctx context.Context, req *server.SwipeRequest) error {
	userId, err := verifyToken(req.Token)
	if err != nil {
		return err
	}
	pa, err := s.repository.GetPendingPairAttemptByUserPair(userId, req.CandidateID)
	if err != nil {
		return err
	}
	var eventType models.PEType
	if req.SwipeVerdict == server.SwipeLike && pa.User1 == userId {
		eventType = models.PETypeUser1Liked
	} else if req.SwipeVerdict == server.SwipeDislike && pa.User1 == userId {
		eventType = models.PETypeUser1Disliked
	} else if req.SwipeVerdict == server.SwipeLike && pa.User2 == userId {
		eventType = models.PETypeUser2Liked
	} else if req.SwipeVerdict == server.SwipeDislike && pa.User2 == userId {
		eventType = models.PETypeUser2Disliked
	}
	err = s.repository.CreateEvent(pa.ID, eventType)
	if err != nil {
		return err
	}

	if req.SwipeVerdict == server.SwipeDislike {
		return s.repository.FinishPairAttempt(pa.ID, models.PAStateMismatch)
	}
	if pa.User1 == userId {
		err = s.notifyLikedUser(ctx, userId, pa.User2)
		if err != nil {
			return err
		}
	} else {
		err = s.notifyMatch(ctx, pa.User1, pa.User2)
		if err != nil {
			return err
		}
		err = s.repository.FinishPairAttempt(pa.ID, models.PAStateMatch)
		if err != nil {
			return err
		}
		err = s.repository.CreateChat(pa.User1, pa.User2)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) notifyLikedUser(ctx context.Context, whoLiked, whomLiked uint64) error {
	prof, err := s.repository.GetProfile(whoLiked)
	if err != nil {
		return err
	}
	photos, err := s.repository.GetUserPhotos(whoLiked)
	if err != nil {
		return err
	}
	link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
	if err != nil {
		return err
	}
	err = s.userInteractor.NotifyLiked(whomLiked, prof.Name, link)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) notifyMatch(ctx context.Context, user1, user2 uint64) error {
	err := s.oneDirectionalNotifyMatch(ctx, user1, user2)
	if err != nil {
		return err
	}
	return s.oneDirectionalNotifyMatch(ctx, user2, user1)
}

func (s *Service) oneDirectionalNotifyMatch(ctx context.Context, sender, receiver uint64) error {
	prof, err := s.repository.GetProfile(sender)
	if err != nil {
		return err
	}
	photos, err := s.repository.GetUserPhotos(sender)
	if err != nil {
		return err
	}
	link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
	if err != nil {
		return err
	}
	err = s.userInteractor.NotifyMatch(receiver, prof.Name, link)
	if err != nil {
		return err
	}
	return nil
}
