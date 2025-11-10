package linkhttp

import "net/http"

func Route(h *LinkCheckerHandler) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /status", h.CheckLinksStatus)

	return router
}
