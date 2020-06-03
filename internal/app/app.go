package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DanielTitkov/netology-slowly/internal/api"
	"github.com/DanielTitkov/netology-slowly/internal/configs"
)

// App holds app configurations
type App struct {
	cfg configs.Config
}

// NewApp created app internal structure and return pointer to App
func NewApp(cfg configs.Config) *App {
	return &App{
		cfg: cfg,
	}
}

// SlowHandler implements "slow" requests logic
func (a *App) SlowHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request api.SlowRequestBody
	err := decoder.Decode(&request)
	if err != nil {
		resp := api.ErrorResponseBody{
			Error: api.ErrorInvalidRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	time.Sleep(time.Duration(request.Timeout) * time.Millisecond)

	resp := api.OkResponseBody{
		Status: api.OkStatus,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
