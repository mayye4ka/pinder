package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	service Service
}

type Service interface {
	RegisterUser(req *RegisterRequest) (*RegisterResponse, error)
	LoginUser(req *LoginRequest) (*LoginResponse, error)

	GetProfile(req *RequestWithToken) (*GetProfileResponse, error)
	UpdProfile(req *UpdProfileRequest) error
	GetPreferences(req *RequestWithToken) (*GetPreferencesResponse, error)
	UpdPreferences(req *UpdPreferencesRequest) error

	NextPartner(req *RequestWithToken) (*NextPartnerResponse, error)
	Swipe(req *SwipeRequest) error
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
		// TODO: good error handling
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
		return s.service.RegisterUser(&req)
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.LoginUser(&req)
	})
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.GetProfile(&req)
	})
}
func (s *Server) updProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdProfileRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.UpdProfile(&req)
	})
}

func (s *Server) getPreferences(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.GetPreferences(&req)
	})
}

func (s *Server) updPreferences(w http.ResponseWriter, r *http.Request) {
	var req UpdPreferencesRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.UpdPreferences(&req)
	})
}

func (s *Server) nextPartner(w http.ResponseWriter, r *http.Request) {
	var req RequestWithToken
	s.commonHandler(w, r, &req, func() (any, error) {
		return s.service.NextPartner(&req)
	})
}

func (s *Server) swipe(w http.ResponseWriter, r *http.Request) {
	var req SwipeRequest
	s.commonHandler(w, r, &req, func() (any, error) {
		return nil, s.service.Swipe(&req)
	})
}

func (s *Server) Start() error {
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/login", s.login)
	// TODO:
	return http.ListenAndServe(":8080", nil)
}
