package server

import (
	"net/http"

	linkhttp "github.com/Negat1v9/link-checker/internal/linkChecker/http"
	"github.com/Negat1v9/link-checker/internal/linkChecker/service"
	"github.com/Negat1v9/link-checker/pkg/logger"
)

func (s *Server) MapHandlers(log *logger.Logger, linkService *service.LinkCheckerService) {
	hander := http.NewServeMux()

	linkHandler := linkhttp.NewLinkCheckerHandler(s.cfg, log, linkService)

	linkRoutes := linkhttp.Route(linkHandler)

	hander.Handle("/links/", http.StripPrefix("/links", linkRoutes))

	s.httpServer.Handler = hander
}
