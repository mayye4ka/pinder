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
		// notify user 2 that he was liked
	} else {
		// mb create chat and notify two users
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
