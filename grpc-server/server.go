package server

import (
	"context"

	public_api "github.com/mayye4ka/pinder-api/api/go"
	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	auth    Authenticator
	service Service
	public_api.UnimplementedPinderServer
}

type Service interface {
	GetProfile(ctx context.Context) (models.ProfileShowcase, error)
	UpdProfile(ctx context.Context, newProfile models.Profile) error
	GetPreferences(ctx context.Context) (models.Preferences, error)
	UpdPreferences(ctx context.Context, newPreferences models.Preferences) error

	AddPhoto(ctx context.Context, photo string) error
	DeletePhoto(ctx context.Context, photoKey string) error
	ReorderPhotos(ctx context.Context, newOrder []string) error

	NextPartner(ctx context.Context) (models.ProfileShowcase, error)
	Swipe(ctx context.Context, candidateId uint64, swipeVerdict models.SwipeVerdict) error

	ListChats(ctx context.Context) ([]models.ChatShowcase, error)
	ListMessages(ctx context.Context, chatId uint64) ([]models.MessageShowcase, error)
	SendMessage(ctx context.Context, chatId uint64, contentType models.MsgContentType, payload string) error
	GetTextFromVoice(ctx context.Context, msgId uint64) (string, bool, error)
}

type Authenticator interface {
	UnpackToken(ctx context.Context, token string) (uint64, error)
	Register(ctx context.Context, phone, password string) (string, error)
	Login(ctx context.Context, phone, password string) (string, error)
}

func (s *Server) Register(ctx context.Context, req *public_api.RegisterRequest) (*public_api.RegisterResponse, error) {
	token, err := s.auth.Register(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.RegisterResponse{
		Token: token,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *public_api.LoginRequest) (*public_api.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.LoginResponse{
		Token: token,
	}, nil
}

func (s *Server) GetUserId(ctx context.Context, req *emptypb.Empty) (*public_api.GetUserIdResponse, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return &public_api.GetUserIdResponse{
		Id: userId,
	}, nil
}

func (s *Server) GetProfile(ctx context.Context, _ *emptypb.Empty) (*public_api.GetProfileResponse, error) {
	profile, err := s.service.GetProfile(ctx)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.GetProfileResponse{
		Profile: profileToProto(profile.Profile),
		Photos:  photosToProto(profile.Photos),
	}, nil
}

func (s *Server) GetPreferences(ctx context.Context, _ *emptypb.Empty) (*public_api.GetPreferencesResponse, error) {
	preferences, err := s.service.GetPreferences(ctx)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.GetPreferencesResponse{
		Preferences: preferencesToProto(preferences),
	}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *public_api.UpdateProfileRequest) (*emptypb.Empty, error) {
	err := s.service.UpdProfile(ctx, protoToProfile(req.NewProfile))
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdatePreferences(ctx context.Context, req *public_api.UpdatePreferencesRequest) (*emptypb.Empty, error) {
	err := s.service.UpdPreferences(ctx, protoToPreferences(req.NewPreferences))
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) AddPhoto(ctx context.Context, req *public_api.AddPhotoRequest) (*emptypb.Empty, error) {
	err := s.service.AddPhoto(ctx, string(req.Photo))
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ReorderPhotos(ctx context.Context, req *public_api.ReorderPhotosRequest) (*emptypb.Empty, error) {
	err := s.service.ReorderPhotos(ctx, req.NewOrder)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeletePhoto(ctx context.Context, req *public_api.DeletePhotoRequest) (*emptypb.Empty, error) {
	err := s.service.DeletePhoto(ctx, req.PhotoKey)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) NextPartner(ctx context.Context, _ *emptypb.Empty) (*public_api.NextPartnerResponse, error) {
	candidate, err := s.service.NextPartner(ctx)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.NextPartnerResponse{
		Candidate: profileShowcaseToCandidate(candidate),
	}, nil
}

func (s *Server) Swipe(ctx context.Context, req *public_api.SwipeRequest) (*emptypb.Empty, error) {
	err := s.service.Swipe(ctx, req.CandidateId, protoToSwipeVerdict(req.SwipeVerdict))
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListChats(ctx context.Context, _ *emptypb.Empty) (*public_api.ListChatsResponse, error) {
	chats, err := s.service.ListChats(ctx)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.ListChatsResponse{
		Chats: chatsToProto(chats),
	}, nil
}

func (s *Server) ListMessages(ctx context.Context, req *public_api.ListMessagesRequest) (*public_api.ListMessagesResponse, error) {
	messages, err := s.service.ListMessages(ctx, req.ChatId)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.ListMessagesResponse{
		Messages: messagesToProto(messages),
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *public_api.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.service.SendMessage(ctx, req.ChatId, protoToMsgContentType(req.ContentType), string(req.Payload))
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetTextFromVoice(ctx context.Context, req *public_api.GetTextFromVoiceRequest) (*public_api.GetTextFromVoiceResponse, error) {
	text, shouldWait, err := s.service.GetTextFromVoice(ctx, req.MessageId)
	if err != nil {
		return nil, errs.ToGrpcError(err)
	}
	return &public_api.GetTextFromVoiceResponse{
		Text:       text,
		ShouldWait: shouldWait,
	}, nil
}
