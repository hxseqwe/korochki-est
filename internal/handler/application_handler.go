package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hxseqwe/korochki-est/internal/model"
	"github.com/hxseqwe/korochki-est/internal/service"
	"net/http"
	"strconv"
)

type ApplicationHandler struct {
	appService   *service.ApplicationService
	sessionStore *sessions.CookieStore
}

func NewApplicationHandler(appService *service.ApplicationService, sessionStore *sessions.CookieStore) *ApplicationHandler {
	return &ApplicationHandler{
		appService:   appService,
		sessionStore: sessionStore,
	}
}

func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var req model.ApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, _ := h.sessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	app, err := h.appService.Create(userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *ApplicationHandler) GetUserApplications(w http.ResponseWriter, r *http.Request) {
	session, _ := h.sessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	apps, err := h.appService.GetUserApplications(userID)
	if err != nil {
		http.Error(w, "Failed to get applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (h *ApplicationHandler) GetAllApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appService.GetAllApplications()
	if err != nil {
		http.Error(w, "Failed to get applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (h *ApplicationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	var req model.StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.appService.UpdateStatus(appID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	var req model.ApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.appService.UpdateApplication(appID, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	if err := h.appService.DeleteApplication(appID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ApplicationHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	var req model.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.appService.AddReview(appID, req.Review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
