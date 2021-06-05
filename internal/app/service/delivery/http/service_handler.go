package http

import (
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

func (sh *ServiceHandler) ServiceClearHandler(w http.ResponseWriter, r *http.Request) {
	errE := sh.serviceUsecase.ClearDatabases()
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Очистка базы успешно завершена"))
}

func (sh *ServiceHandler) ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, errE := sh.serviceUsecase.GetServiceStatus()
	if errE != nil {
		w.WriteHeader(errE.HttpError)
		w.Write(errors.JSONError(errE))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(errors.JSONSuccess("Кол-во записей в базе данных, включая помеченные как \"удалённые\".", status))
}