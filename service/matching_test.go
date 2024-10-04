package service

import (
	"github.com/mayye4ka/pinder/models"
)

func (s *ServiceTestSuite) TestNextPartner_ReturnsHangingPartner() {
	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return([]string{photo1, photo2}, nil)
	for _, k := range photos {
		s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, k.Key).Return(k.Link, nil)
	}
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         userName,
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID:       user2Id,
			Name:         userName,
			Gender:       models.GenderMale,
			Age:          20,
			Bio:          "bio",
			LocationLat:  123,
			LocationLon:  456,
			LocationName: "Kolbasino neighbourghood",
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsWhoLikedMe() {
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(nil, nil)

	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return([]string{photo1, photo2}, nil)
	for _, k := range photos {
		s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, k.Key).Return(k.Link, nil)
	}
	s.repoMock.EXPECT().GetLatestPairAttempt(user2Id, userId).Return(models.PairAttempt{ID: 1}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeSentToUser2).Return(nil)
	s.repoMock.EXPECT().GetWhoLikedMe(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         userName,
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID:       user2Id,
			Name:         userName,
			Gender:       models.GenderMale,
			Age:          20,
			Bio:          "bio",
			LocationLat:  123,
			LocationLon:  456,
			LocationName: "Kolbasino neighbourghood",
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsNewPair() {
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(nil, nil)
	s.repoMock.EXPECT().GetWhoLikedMe(userId).Return(nil, nil)

	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return([]string{photo1, photo2}, nil)
	for _, k := range photos {
		s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, k.Key).Return(k.Link, nil)
	}
	s.repoMock.EXPECT().ChooseCandidateAndCreatePairAttempt(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         userName,
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID:       user2Id,
			Name:         userName,
			Gender:       models.GenderMale,
			Age:          20,
			Bio:          "bio",
			LocationLat:  123,
			LocationLon:  456,
			LocationName: "Kolbasino neighbourghood",
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestSwipe_First_Like() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: userId, User2: user2Id}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser1Liked).Return(nil)
	s.repoMock.EXPECT().GetProfile(userId).Return(models.Profile{UserID: userId, Name: userName}, nil)
	s.repoMock.EXPECT().GetUserPhotos(userId).Return([]string{photo1, photo2}, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.userNotifierMock.EXPECT().NotifyLiked(user2Id, models.LikeNotification{
		Name:  userName,
		Photo: photo1Link,
	}).Return(nil)

	err := s.service.Swipe(user1Ctx, user2Id, models.SwipeVerdictLike)

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_First_Dislike() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: userId, User2: user2Id}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser1Disliked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMismatch).Return(nil)

	err := s.service.Swipe(user1Ctx, user2Id, models.SwipeVerdictDislike)

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_Second_Like() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: user2Id, User2: userId}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser2Liked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMatch).Return(nil)
	s.repoMock.EXPECT().CreateChat(user2Id, userId).Return(nil)

	s.repoMock.EXPECT().GetProfile(userId).Return(models.Profile{UserID: userId, Name: userName}, nil)
	s.repoMock.EXPECT().GetUserPhotos(userId).Return([]string{photo1, photo2}, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.userNotifierMock.EXPECT().NotifyMatch(userId, models.MatchNotification{
		Name:  userName,
		Photo: photo1Link,
	}).Return(nil)

	s.repoMock.EXPECT().GetProfile(user2Id).Return(models.Profile{UserID: user2Id, Name: userName}, nil)
	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return([]string{photo1, photo2}, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.userNotifierMock.EXPECT().NotifyMatch(user2Id, models.MatchNotification{
		Name:  userName,
		Photo: photo1Link,
	}).Return(nil)

	err := s.service.Swipe(user1Ctx, user2Id, models.SwipeVerdictLike)

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_Second_Dislike() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: user2Id, User2: userId}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser2Disliked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMismatch).Return(nil)

	err := s.service.Swipe(user1Ctx, user2Id, models.SwipeVerdictDislike)

	s.Nil(err)
}
