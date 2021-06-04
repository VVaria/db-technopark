package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/internal/app/service"
	"github.com/VVaria/db-technopark/internal/app/tools/errors"
)

type ServiceHandler struct {
	serviceUsecase service.ServiceUsecase
}

func NewServiceHandler(serviceUsecase service.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUsecase: serviceUsecase,
	}
}

func (sh *ServiceHandler) Configure(r *mux.Router) {
	r.HandleFunc("/service/clear", sh.ServiceClearHandler).Methods(http.MethodPost)
	r.HandleFunc("/service/status", sh.ServiceStatusHandler).Methods(http.MethodGet)
}
