package service

import (
	"pinder/models"
	"pinder/server"
)

func (s *ServiceTestSuite) TestNextPartner_ReturnsHandingPartner() {
	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return(photoKeys, nil)
	for i, k := range photoKeys {
		s.fsMock.EXPECT().MakeLink(testCtx, k).Return(photoLinks[i], nil)
	}
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         "name",
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	resp, err := s.service.NextPartner(testCtx, &server.RequestWithToken{
		Token: token,
	})

	s.Nil(err)
	s.Equal(server.Profile{
		Name:         "name",
		Gender:       server.GenderMale,
		Age:          20,
		Bio:          "bio",
		Photos:       photoLinks,
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, resp.Partner)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsWhoLikedMe() {
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(nil, nil)

	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return(photoKeys, nil)
	for i, k := range photoKeys {
		s.fsMock.EXPECT().MakeLink(testCtx, k).Return(photoLinks[i], nil)
	}
	s.repoMock.EXPECT().GetLatestPairAttempt(user2Id, userId).Return(models.PairAttempt{ID: 1}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeSentToUser2).Return(nil)
	s.repoMock.EXPECT().GetWhoLikedMe(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         "name",
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	resp, err := s.service.NextPartner(testCtx, &server.RequestWithToken{
		Token: token,
	})

	s.Nil(err)
	s.Equal(server.Profile{
		Name:         "name",
		Gender:       server.GenderMale,
		Age:          20,
		Bio:          "bio",
		Photos:       photoLinks,
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, resp.Partner)
}

func (s *ServiceTestSuite) TestNextPartner_ReturnsNewPair() {
	s.repoMock.EXPECT().GetHangingPartner(userId).Return(nil, nil)
	s.repoMock.EXPECT().GetWhoLikedMe(userId).Return(nil, nil)

	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return(photoKeys, nil)
	for i, k := range photoKeys {
		s.fsMock.EXPECT().MakeLink(testCtx, k).Return(photoLinks[i], nil)
	}
	s.repoMock.EXPECT().ChooseCandidateAndCreatePairAttempt(userId).Return(&models.Profile{
		UserID:       user2Id,
		Name:         "name",
		Gender:       models.GenderMale,
		Age:          20,
		Bio:          "bio",
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, nil)

	resp, err := s.service.NextPartner(testCtx, &server.RequestWithToken{
		Token: token,
	})

	s.Nil(err)
	s.Equal(server.Profile{
		Name:         "name",
		Gender:       server.GenderMale,
		Age:          20,
		Bio:          "bio",
		Photos:       photoLinks,
		LocationLat:  123,
		LocationLon:  456,
		LocationName: "Kolbasino neighbourghood",
	}, resp.Partner)
}

func (s *ServiceTestSuite) TestSwipe_First_Like() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: userId, User2: user2Id}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser1Liked).Return(nil)

	err := s.service.Swipe(testCtx, &server.SwipeRequest{
		Token:        token,
		CandidateID:  user2Id,
		SwipeVerdict: server.SwipeLike,
	})

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_First_Dislike() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: userId, User2: user2Id}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser1Disliked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMismatch).Return(nil)

	err := s.service.Swipe(testCtx, &server.SwipeRequest{
		Token:        token,
		CandidateID:  user2Id,
		SwipeVerdict: server.SwipeDislike,
	})

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_Second_Like() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: user2Id, User2: userId}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser2Liked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMatch).Return(nil)
	s.repoMock.EXPECT().CreateChat(user2Id, userId).Return(nil)

	err := s.service.Swipe(testCtx, &server.SwipeRequest{
		Token:        token,
		CandidateID:  user2Id,
		SwipeVerdict: server.SwipeLike,
	})

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSwipe_Second_Dislike() {
	s.repoMock.EXPECT().GetPendingPairAttemptByUserPair(userId, user2Id).Return(models.PairAttempt{ID: 1, User1: user2Id, User2: userId}, nil)
	s.repoMock.EXPECT().CreateEvent(uint64(1), models.PETypeUser2Disliked).Return(nil)
	s.repoMock.EXPECT().FinishPairAttempt(uint64(1), models.PAStateMismatch).Return(nil)

	err := s.service.Swipe(testCtx, &server.SwipeRequest{
		Token:        token,
		CandidateID:  user2Id,
		SwipeVerdict: server.SwipeDislike,
	})

	s.Nil(err)
}
