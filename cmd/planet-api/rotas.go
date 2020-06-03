package main

import (
	"desafio-b2w/core"
	"desafio-b2w/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func strToUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 64)
	return uint(n)
}

func rotas(f *core.Facade) chi.Router {
	r := chi.NewRouter()
	r.Post("/planets", addPlanetHandler(f))
	r.Get("/planets", allPlanetsHandler(f))
	r.Get("/planets/{id}", planetByIDHandler(f))
	r.Delete("/planets/{id}", removePlanetHandler(f))
	return r
}

func allPlanetsHandler(f *core.Facade) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			result interface{}
			err    error
		)
		name := r.URL.Query().Get("name")
		if name == "" {
			result, err = f.Planets().List()
		} else {
			result, err = f.Planets().GetByName(name)
		}
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonWrite(w, http.StatusOK, result)
	}
}

func addPlanetHandler(f *core.Facade) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var planet models.Planet
		if err := jsonRead(w, r, &planet); err != nil {
			return
		}
		id, err := f.Planets().Add(planet)
		if err != nil {
			jsonError(w, err)
			return
		}
		planet, err = f.Planets().Get(id.Hex())
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonWrite(w, http.StatusCreated, planet)
	}
}

func planetByIDHandler(f *core.Facade) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := f.Planets().Get(chi.URLParam(r, "id"))
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonWrite(w, http.StatusOK, p)
	}
}

func removePlanetHandler(f *core.Facade) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f.Planets().Remove(chi.URLParam(r, "id"))
		if err != nil {
			jsonError(w, err)
		}
	}
}
