package main

import (
	"desafio-b2w/errors"
	"encoding/json"
	"log"
	"net/http"
)

// ResponseError represents error response.
type ResponseError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (re ResponseError) Error() string {
	return re.Message
}

func jsonWrite(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	jsonEnc := json.NewEncoder(w)
	err := jsonEnc.Encode(v)
	if err != nil {
		log.Println("Erro ao codificar JSON:", err)
	}
}

func jsonWriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	e := ResponseError{status, msg}
	jsonEnc := json.NewEncoder(w)
	err := jsonEnc.Encode(e)
	if err != nil {
		log.Println("Erro ao codificar JSON:", err)
	}
}

func jsonError(w http.ResponseWriter, e error) {
	var re ResponseError
	switch err := e.(type) {
	default:
		re.Status = http.StatusInternalServerError
		re.Message = e.Error()
	case errors.InputError:
		re.Status = http.StatusBadRequest
		re.Message = err.Error()
	case errors.NotFoundError:
		re.Status = http.StatusNotFound
		re.Message = err.Error()
	case errors.ConflictError:
		re.Status = http.StatusConflict
		re.Message = err.Error()
	case errors.InternalError:
		re.Status = http.StatusInternalServerError
		re.Message = err.Error()
	case ResponseError:
		re = err
	}

	jsonWriteError(w, re.Status, re.Message)
}

func jsonRead(w http.ResponseWriter, r *http.Request, v interface{}) (err error) {
	jsonDec := json.NewDecoder(r.Body)
	err = jsonDec.Decode(&v)
	if err != nil {
		jsonWriteError(w, http.StatusBadRequest, "Corpo de requisição inválido.")
	}
	return
}
