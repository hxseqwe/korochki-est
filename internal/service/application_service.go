package service

import (
	"errors"
	"github.com/hxseqwe/korochki-est/internal/model"
	"github.com/hxseqwe/korochki-est/internal/repository"
	"time"
)

type ApplicationService struct {
	appRepo *repository.ApplicationRepository
}

func NewApplicationService(appRepo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		appRepo: appRepo,
	}
}

func (s *ApplicationService) Create(userID int, req *model.ApplicationRequest) (*model.Application, error) {
	if req.CourseName == "" {
		return nil, errors.New("course name is required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	if req.PaymentMethod != "cash" && req.PaymentMethod != "transfer" {
		return nil, errors.New("invalid payment method")
	}

	app := &model.Application{
		UserID:        userID,
		CourseName:    req.CourseName,
		StartDate:     startDate,
		PaymentMethod: req.PaymentMethod,
		Status:        "new",
	}

	err = s.appRepo.Create(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *ApplicationService) GetUserApplications(userID int) ([]*model.Application, error) {
	return s.appRepo.GetByUserID(userID)
}

func (s *ApplicationService) GetAllApplications() ([]*model.Application, error) {
	return s.appRepo.GetAll()
}

func (s *ApplicationService) UpdateStatus(appID int, status string) error {
	validStatuses := map[string]bool{
		"new": true, "in_progress": true, "completed": true, "rejected": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	return s.appRepo.UpdateStatus(appID, status)
}

func (s *ApplicationService) UpdateApplication(appID int, req *model.ApplicationRequest) error {
	if req.CourseName == "" {
		return errors.New("course name is required")
	}

	if req.PaymentMethod != "cash" && req.PaymentMethod != "transfer" {
		return errors.New("invalid payment method")
	}

	return s.appRepo.UpdateApplication(appID, req.CourseName, req.StartDate, req.PaymentMethod)
}

func (s *ApplicationService) DeleteApplication(appID int) error {
	return s.appRepo.DeleteApplication(appID)
}

func (s *ApplicationService) AddReview(appID int, review string) error {
	if review == "" {
		return errors.New("review cannot be empty")
	}

	app, err := s.appRepo.GetByID(appID)
	if err != nil {
		return err
	}

	if app.Status != "completed" {
		return errors.New("can only review completed courses")
	}

	return s.appRepo.AddReview(appID, review)
}
