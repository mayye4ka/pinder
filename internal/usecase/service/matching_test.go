package service

import "github.com/mayye4ka/pinder/internal/models"

var (
	PAID = uint64(777)
)

func (s *ServiceTestSuite) TestNextPartner_InvalidUser() {
	s.repoMock.EXPECT().GetProfile(user1Ctx, userId).Return(models.Profile{}, nil)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, userId).Return(models.Preferences{}, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Equal("incomplete profile", err.Error())
	s.Equal(models.ProfileShowcase{}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsHangingPartner() {
	s.repoMock.EXPECT().GetProfile(user1Ctx, userId).Return(models.Profile{
		UserID: userId,
	}, nil)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, userId).Return(models.Preferences{
		UserID: userId,
	}, nil)

	s.repoMock.EXPECT().GetPendingPairAttempts(user1Ctx, userId).Return([]models.PairAttempt{{
		ID:    PAID,
		User2: user2Id,
	}}, nil)
	s.repoMock.EXPECT().GetLastEvent(user1Ctx, PAID).Return(models.PairEvent{EventType: models.PETypeSentToUser1}, nil)

	s.repoMock.EXPECT().GetProfile(user1Ctx, user2Id).Return(models.Profile{UserID: user2Id}, nil)
	s.repoMock.EXPECT().GetUserPhotos(user1Ctx, user2Id).Return([]string{photo1, photo2}, nil)

	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo2).Return(photo2Link, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID: user2Id,
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsWhoLikedMe() {
	s.repoMock.EXPECT().GetProfile(user1Ctx, userId).Return(models.Profile{
		UserID: userId,
	}, nil)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, userId).Return(models.Preferences{
		UserID: userId,
	}, nil)

	s.repoMock.EXPECT().GetPendingPairAttempts(user1Ctx, userId).Return([]models.PairAttempt{}, nil)

	s.repoMock.EXPECT().GetWhoLikedMe(user1Ctx, userId).Return(user2Id, nil)
	s.repoMock.EXPECT().GetLatestPairAttempt(user1Ctx, user2Id, userId).Return(models.PairAttempt{ID: PAID}, nil)
	s.repoMock.EXPECT().CreateEvent(user1Ctx, PAID, models.PETypeSentToUser2).Return(nil)

	s.repoMock.EXPECT().GetProfile(user1Ctx, user2Id).Return(models.Profile{UserID: user2Id}, nil)
	s.repoMock.EXPECT().GetUserPhotos(user1Ctx, user2Id).Return([]string{photo1, photo2}, nil)

	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo2).Return(photo2Link, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID: user2Id,
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsNewPair() {
	s.repoMock.EXPECT().GetProfile(user1Ctx, userId).Return(models.Profile{
		UserID: userId,
		Age:    19,
	}, nil).Times(2)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, userId).Return(models.Preferences{
		UserID: userId,
		MinAge: 18,
		MaxAge: 20,
	}, nil).Times(2)

	s.repoMock.EXPECT().GetPendingPairAttempts(user1Ctx, userId).Return([]models.PairAttempt{}, nil)
	s.repoMock.EXPECT().GetWhoLikedMe(user1Ctx, userId).Return(uint64(0), nil)

	s.repoMock.EXPECT().GetAllValidUsers(user1Ctx).Return([]uint64{user2Id}, nil)
	s.repoMock.EXPECT().GetProfile(user1Ctx, user2Id).Return(models.Profile{
		UserID: user2Id,
		Age:    19,
	}, nil).Times(2)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, user2Id).Return(models.Preferences{
		UserID: user2Id,
		MinAge: 18,
		MaxAge: 20,
	}, nil)

	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(user1Ctx, userId, user2Id).Return(models.PairAttempt{}, nil)
	s.repoMock.EXPECT().GetLatestPairAttemptByUserPair(user1Ctx, userId, user2Id).Return(models.PairAttempt{}, nil)
	s.repoMock.EXPECT().CreatePairAttempt(user1Ctx, userId, user2Id).Return(models.PairAttempt{ID: PAID}, nil)
	s.repoMock.EXPECT().CreateEvent(user1Ctx, PAID, models.PETypePACreated).Return(nil)
	s.repoMock.EXPECT().CreateEvent(user1Ctx, PAID, models.PETypeSentToUser1).Return(nil)

	s.repoMock.EXPECT().GetUserPhotos(user1Ctx, user2Id).Return([]string{photo1, photo2}, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo2).Return(photo2Link, nil)

	candidate, err := s.service.NextPartner(user1Ctx)

	s.Nil(err)
	s.Equal(models.ProfileShowcase{
		Profile: models.Profile{
			UserID: user2Id,
			Age:    19,
		},
		Photos: photos,
	}, candidate)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsValuableAdvice() {
	s.repoMock.EXPECT().GetProfile(user1Ctx, userId).Return(models.Profile{
		UserID: userId,
	}, nil).Times(2)
	s.repoMock.EXPECT().GetPreferences(user1Ctx, userId).Return(models.Preferences{
		UserID: userId,
	}, nil).Times(2)

	s.repoMock.EXPECT().GetPendingPairAttempts(user1Ctx, userId).Return([]models.PairAttempt{}, nil)
	s.repoMock.EXPECT().GetWhoLikedMe(user1Ctx, userId).Return(uint64(0), nil)

	s.repoMock.EXPECT().GetAllValidUsers(user1Ctx).Return([]uint64{}, nil)

	candidate, err := s.service.NextPartner(user1Ctx)
	s.Equal("lower your expectations to zero", err.Error())
	s.Equal(models.ProfileShowcase{}, candidate)
}
