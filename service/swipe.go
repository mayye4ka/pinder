package service

import (
	"github.com/mayye4ka/pinder/models"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (s *Service) Swipe(ctx context.Context, candidateId uint64, swipeVerdict models.SwipeVerdict) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	pa, err := s.repository.GetPendingPairAttemptByUserPair(userId, candidateId)
	if err != nil {
		return errors.Wrap(err, "can't get pending pa by user pair")
	}
	var eventType models.PEType
	if swipeVerdict == models.SwipeVerdictLike && pa.User1 == userId {
		eventType = models.PETypeUser1Liked
	} else if swipeVerdict == models.SwipeVerdictDislike && pa.User1 == userId {
		eventType = models.PETypeUser1Disliked
	} else if swipeVerdict == models.SwipeVerdictLike && pa.User2 == userId {
		eventType = models.PETypeUser2Liked
	} else if swipeVerdict == models.SwipeVerdictDislike && pa.User2 == userId {
		eventType = models.PETypeUser2Disliked
	}
	err = s.repository.CreateEvent(pa.ID, eventType)
	if err != nil {
		return errors.Wrap(err, "can't create event")
	}

	if swipeVerdict == models.SwipeVerdictDislike {
		err := s.repository.FinishPairAttempt(pa.ID, models.PAStateMismatch)
		if err != nil {
			return errors.Wrap(err, "can't finish pair attempt")
		}
		return nil
	}
	if pa.User1 == userId {
		err = s.notifyLikedUser(ctx, userId, pa.User2)
		if err != nil {
			return errors.Wrap(err, "can't notify liked user")
		}
	} else {
		err = s.notifyMatch(ctx, pa.User1, pa.User2)
		if err != nil {
			return errors.Wrap(err, "can't notify match")
		}
		err = s.repository.FinishPairAttempt(pa.ID, models.PAStateMatch)
		if err != nil {
			return errors.Wrap(err, "can't finish pair attempt")
		}
		err = s.repository.CreateChat(pa.User1, pa.User2)
		if err != nil {
			return errors.Wrap(err, "can't create chat")
		}
	}
	return nil
}

func (s *Service) notifyLikedUser(ctx context.Context, whoLiked, whomLiked uint64) error {
	prof, err := s.repository.GetProfile(whoLiked)
	if err != nil {
		return errors.Wrap(err, "can't get profile")
	}
	photos, err := s.repository.GetUserPhotos(whoLiked)
	if err != nil {
		return errors.Wrap(err, "can't get user photos")
	}
	link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
	if err != nil {
		return errors.Wrap(err, "can't make profile photo link")
	}
	err = s.userNotifier.NotifyLiked(whomLiked, models.LikeNotification{
		Name:  prof.Name,
		Photo: link,
	})
	if err != nil {
		return errors.Wrap(err, "can't notify liked user")
	}
	return nil
}

func (s *Service) notifyMatch(ctx context.Context, user1, user2 uint64) error {
	err := s.oneDirectionalNotifyMatch(ctx, user1, user2)
	if err != nil {
		return errors.Wrap(err, "can't send match notification")
	}
	err = s.oneDirectionalNotifyMatch(ctx, user2, user1)
	if err != nil {
		return errors.Wrap(err, "can't send match notification")
	}
	return nil
}

func (s *Service) oneDirectionalNotifyMatch(ctx context.Context, sender, receiver uint64) error {
	prof, err := s.repository.GetProfile(sender)
	if err != nil {
		return errors.Wrap(err, "can't get profile")
	}
	photos, err := s.repository.GetUserPhotos(sender)
	if err != nil {
		return errors.Wrap(err, "can't get user photos")
	}
	link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
	if err != nil {
		return errors.Wrap(err, "can't make profile photo link")
	}
	err = s.userNotifier.NotifyMatch(receiver, models.MatchNotification{
		Name:  prof.Name,
		Photo: link,
	})
	if err != nil {
		return errors.Wrap(err, "can't notify match")
	}
	return nil
}
