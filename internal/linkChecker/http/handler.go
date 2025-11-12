package linkhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Negat1v9/link-checker/config"
	linkmodel "github.com/Negat1v9/link-checker/internal/linkChecker/model"
	"github.com/Negat1v9/link-checker/internal/linkChecker/service"
	"github.com/Negat1v9/link-checker/pkg/logger"
)

type LinkCheckerHandler struct {
	cfg         *config.ServerCfg
	log         *logger.Logger
	linkService *service.LinkCheckerService
}

func NewLinkCheckerHandler(cfg *config.ServerCfg, log *logger.Logger, service *service.LinkCheckerService) *LinkCheckerHandler {
	return &LinkCheckerHandler{
		cfg:         cfg,
		log:         log,
		linkService: service,
	}
}

func (h *LinkCheckerHandler) CheckLinksStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*h.cfg.DefaultWriteTimeOut)
	defer cancel()

	var req linkmodel.CheckLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debugf("linkhttp.CheckLinksStatus: %v", err)
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	result := h.linkService.CheckLinks(ctx, req.Links)

	json.NewEncoder(w).Encode(result)
}

func (h *LinkCheckerHandler) CreateLinksGroupReport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*h.cfg.DefaultWriteTimeOut)
	defer cancel()

	var req linkmodel.CheckLinkListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debugf("linkhttp.CreateLinksGroupReport: %v", err)
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	pdf, err := h.linkService.CreateLinksGroupPdfReport(ctx, req.LinksList)
	if err != nil {
		h.log.Errorf("linkhttp.CreateLinksGroupReport: %v", err)
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "Content-Disposition: attachment; filename=\"report.pdf\"")

	if err = pdf.Output(w); err != nil {
		h.log.Errorf("linkhttp.CreateLinksGroupReport: %v", err)
		w.WriteHeader(500)
	}

}
