package linkhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Negat1v9/link-checker/config"
	linkmodel "github.com/Negat1v9/link-checker/internal/linkChecker/model"
	"github.com/Negat1v9/link-checker/internal/linkChecker/service"
)

type LinkCheckerHandler struct {
	cfg         *config.ServerCfg
	linkService *service.LinkCheckerService
}

func NewLinkCheckerHandler(cfg *config.ServerCfg, service *service.LinkCheckerService) *LinkCheckerHandler {
	return &LinkCheckerHandler{
		cfg:         cfg,
		linkService: service,
	}
}

func (h *LinkCheckerHandler) CheckLinksStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*h.cfg.DefaultWriteTimeOut)
	defer cancel()

	var req linkmodel.CheckLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	result := h.linkService.CheckLinks(ctx, req.Links)

	json.NewEncoder(w).Encode(result)
}
