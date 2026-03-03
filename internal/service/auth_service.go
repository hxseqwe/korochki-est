package service

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/hxseqwe/korochki-est/internal/model"
	"github.com/hxseqwe/korochki-est/internal/repository"
	"net/http"
)

type AuthService struct {
	userRepo *repository.UserRepository
	store    *sessions.CookieStore
}

func NewAuthService(userRepo *repository.UserRepository, store *sessions.CookieStore) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		store:    store,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	if len(req.Login) < 6 {
		return nil, errors.New("login must be at least 6 characters")
	}

	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	exists, err := s.userRepo.IsLoginExists(req.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("login already exists")
	}

	user := &model.User{
		Login:    req.Login,
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		IsAdmin:  false,
	}

	err = s.userRepo.Create(user, req.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := s.userRepo.FindByLogin(req.Login)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.userRepo.ValidatePassword(user, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) SetSession(w http.ResponseWriter, r *http.Request, user *model.User) error {
	session, _ := s.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["is_admin"] = user.IsAdmin
	return session.Save(r, w)
}

func (s *AuthService) ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.store.Get(r, "session")
	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}

func (s *AuthService) GetCurrentUser(r *http.Request) (*model.User, error) {
	session, _ := s.store.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return nil, errors.New("not authenticated")
	}

	return s.userRepo.FindByID(userID)
}

func (s *AuthService) IsAdmin(r *http.Request) bool {
	session, _ := s.store.Get(r, "session")
	isAdmin, ok := session.Values["is_admin"].(bool)
	return ok && isAdmin
}
