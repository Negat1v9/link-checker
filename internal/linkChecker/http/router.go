package linkhttp

import "net/http"

func Route(h *LinkCheckerHandler) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /status", h.CheckLinksStatus)
	router.HandleFunc("GET /group/report", h.CreateLinksGroupReport)

	return router
}
