package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/timickb/txclient/internal/domain"
	"github.com/timickb/txclient/internal/interfaces"
	"io"
	"net/http"
)

type RequestService interface {
	io.Closer
	SendRequest(ctx context.Context, r domain.ClientRequest) error
}

type Server struct {
	service RequestService
	logger  interfaces.Logger
	ctx     context.Context
}

func New(ctx context.Context, worker RequestService, logger interfaces.Logger) *Server {
	return &Server{
		service: worker,
		logger:  logger,
		ctx:     ctx,
	}
}

func (s *Server) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/transaction", s.transaction)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	return s.service.Close()
}

func (s *Server) transaction(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("POST /transaction")

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	var request domain.ClientRequest

	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(domain.HttpResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	s.logger.Info("Request parsed: ", request)

	if err := s.service.SendRequest(s.ctx, request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(domain.HttpResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if err := encoder.Encode(domain.HttpResponse{
		Status:  http.StatusOK,
		Message: "ok",
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
