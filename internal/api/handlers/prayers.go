package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/catdevman/prayer-journal/internal/models"
	"github.com/catdevman/prayer-journal/internal/repository"
	"github.com/google/uuid"
)

type PrayerHandler struct {
	Repo repository.PrayerRepository
}

func NewPrayerHandler(repo repository.PrayerRepository) *PrayerHandler {
	return &PrayerHandler{Repo: repo}
}

// CreatePrayer handles POST /prayers
func (h *PrayerHandler) CreatePrayer(w http.ResponseWriter, r *http.Request) {
	// 1. Decode Body
	var req models.Prayer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Enrich Model (Business Logic)
	// Get UserID from Context (assumes Middleware set this)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		// Fallback for dev/local if auth isn't running
		userID = "anonymous"
		// In production, return 401 Unauthorized
	}

	req.ID = uuid.NewString()
	req.UserID = userID
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if req.Status == "" {
		req.Status = models.StatusActive
	}

	// 3. Save
	if err := h.Repo.SavePrayer(r.Context(), &req); err != nil {
		http.Error(w, "Failed to save prayer", http.StatusInternalServerError)
		return
	}

	// 4. Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// ListPrayers handles GET /prayers
func (h *PrayerHandler) ListPrayers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		userID = "anonymous"
	}

	prayers, err := h.Repo.GetPrayersByUser(r.Context(), userID, 50)
	if err != nil {
		http.Error(w, "Failed to fetch prayers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prayers)
}
