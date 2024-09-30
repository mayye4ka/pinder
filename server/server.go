package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	service        Service
	userInteractor UserInteractor
}

type Service interface {
	RegisterUser(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	LoginUser(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetUserId(ctx context.Context, req *RequestWithToken) (*GetUserIdResponse, error)

	GetProfile(ctx context.Context, req *RequestWithToken) (*GetProfileResponse, error)
	UpdProfile(ctx context.Context, req *UpdProfileRequest) error
	GetPreferences(ctx context.Context, req *RequestWithToken) (*GetPreferencesResponse, error)
	UpdPreferences(ctx context.Context, req *UpdPreferencesRequest) error

	AddPhoto(ctx context.Context, req *AddPhotoRequest) error
	DeletePhoto(ctx context.Context, req *DelPhotoRequest) error

	NextPartner(ctx context.Context, req *RequestWithToken) (*NextPartnerResponse, error)
	Swipe(ctx context.Context, req *SwipeRequest) error

	ListChats(ctx context.Context, req *RequestWithToken) (*ListChatsResponse, error)
	ListMessages(ctx context.Context, req *ListMessagesRequest) (*ListMessagesResponse, error)
	SendMessage(ctx context.Context, req *SendMessageRequest) error
}

type UserInteractor interface {
	AddUser(userId uint64, conn *websocket.Conn)
}

func New(svc Service, userInteractor UserInteractor) *Server {
	return &Server{
		service:        svc,
		userInteractor: userInteractor,
	}
}

func (s *Server) commonHandler(w http.ResponseWriter, r *http.Request, reqPtr any, serviceMethod func() (any, error)) {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(reqPtr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := serviceMethod()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if resp == nil {
		resp = struct{}{}
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.RegisterUser(r.Context(), &req)
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.LoginUser(r.Context(), &req)
	})
}

func (s *Server) getUserId(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.GetUserId(r.Context(), &req)
	})
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.GetProfile(r.Context(), &req)
	})
}
func (s *Server) updProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdProfileRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.UpdProfile(r.Context(), &req)
	})
}

func (s *Server) getPreferences(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.GetPreferences(r.Context(), &req)
	})
}

func (s *Server) updPreferences(w http.ResponseWriter, r *http.Request) {
	var req UpdPreferencesRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.UpdPreferences(r.Context(), &req)
	})
}

func (s *Server) nextPartner(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.NextPartner(r.Context(), &req)
	})
}

func (s *Server) swipe(w http.ResponseWriter, r *http.Request) {
	var req SwipeRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.Swipe(r.Context(), &req)
	})
}

func (s *Server) addPhoto(w http.ResponseWriter, r *http.Request) {
	var req AddPhotoRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.AddPhoto(r.Context(), &req)
	})
}

func (s *Server) deletePhoto(w http.ResponseWriter, r *http.Request) {
	var req DelPhotoRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.DeletePhoto(r.Context(), &req)
	})
}

func (s *Server) listChats(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.ListChats(r.Context(), &req)
	})
}

func (s *Server) listMessages(w http.ResponseWriter, r *http.Request) {
	var req ListMessagesRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.ListMessages(r.Context(), &req)
	})
}

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.SendMessage(r.Context(), &req)
	})
}

var upgrader = websocket.Upgrader{}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := s.service.GetUserId(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	s.userInteractor.AddUser(resp.ID, conn)
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func (s *Server) Start() error {
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/login", s.login)
	http.HandleFunc("/get_user_id", s.getUserId)

	http.HandleFunc("/get_profile", s.getProfile)
	http.HandleFunc("/get_preferences", s.getPreferences)
	http.HandleFunc("/upd_profile", s.updProfile)
	http.HandleFunc("/upd_preferences", s.updPreferences)
	http.HandleFunc("/upd_photo", s.addPhoto)
	http.HandleFunc("/del_photo", s.deletePhoto)

	http.HandleFunc("/next_partner", s.nextPartner)
	http.HandleFunc("/swipe", s.swipe)

	http.HandleFunc("/list_chats", s.listChats)
	http.HandleFunc("/list_messages", s.listMessages)
	http.HandleFunc("/send_message", s.sendMessage)

	http.HandleFunc("/websocket", s.websocket)

	http.HandleFunc("/", s.hello)
	return http.ListenAndServe(":8080", nil)
}
