package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	service Service
}

type Service interface {
	RegisterUser(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	LoginUser(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	GetProfile(ctx context.Context, req *RequestWithToken) (*GetProfileResponse, error)
	UpdProfile(ctx context.Context, req *UpdProfileRequest) error
	GetPreferences(ctx context.Context, req *RequestWithToken) (*GetPreferencesResponse, error)
	UpdPreferences(ctx context.Context, req *UpdPreferencesRequest) error
	UpdateProfilePhoto(ctx context.Context, token string, photo []byte) error
	DeleteProfilePhoto(ctx context.Context, token string) error

	NextPartner(ctx context.Context, req *RequestWithToken) (*NextPartnerResponse, error)
	Swipe(ctx context.Context, req *SwipeRequest) error
}

func New(svc Service) *Server {
	return &Server{
		service: svc,
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

func (s *Server) updateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if len(r.Form["photo"]) != 1 || len(r.Form["photo"][0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad/no photo provided")
		return
	}
	photoBytes := []byte(r.Form["photo"][0])
	token := r.Header.Get("token")
	err = s.service.UpdateProfilePhoto(r.Context(), token, photoBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (s *Server) deleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	err := s.service.DeleteProfilePhoto(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func (s *Server) Start() error {
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/login", s.login)

	http.HandleFunc("/get_profile", s.getProfile)
	http.HandleFunc("/get_preferences", s.getPreferences)
	http.HandleFunc("/upd_profile", s.updProfile)
	http.HandleFunc("/upd_preferences", s.updPreferences)
	http.HandleFunc("/upd_photo", s.updateProfilePhoto)
	http.HandleFunc("/del_photo", s.deleteProfilePhoto)

	http.HandleFunc("/next_partner", s.nextPartner)
	http.HandleFunc("/swipe", s.swipe)
	http.HandleFunc("/", s.hello)
	return http.ListenAndServe(":8080", nil)
}
