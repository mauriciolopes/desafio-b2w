package core

import (
	"desafio-b2w/db"
	cerrs "desafio-b2w/errors"
	"desafio-b2w/models"
	"desafio-b2w/validate"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrPlanetConflict is planet's name conflict error.
	ErrPlanetConflict = cerrs.ConflictError("Já existe um planeta com o mesmo nome")

	// ErrPlanetNotFound is planet not found error.
	ErrPlanetNotFound = cerrs.NotFoundError("Planeta não encontrado")
)

// PlanetService represents a planet service.
type PlanetService struct {
	f *Facade
}

// Add adds a planet.
func (ps PlanetService) Add(p models.Planet) (primitive.ObjectID, error) {
	p.Name = strings.TrimSpace(p.Name)
	if err := validate.Exec(p); err != nil {
		return primitive.NilObjectID, err
	}
	plRepo := db.NewPlanet(ps.f.db.Get())
	existent, err := plRepo.FindByName(p.Name)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return primitive.NilObjectID, cerrs.UnexpectedError(err)
	}
	if !existent.ID.IsZero() {
		return primitive.NilObjectID, ErrPlanetConflict
	}

	fa, err := ps.getFilmAppearances(p.Name)
	if err != nil {
		return primitive.NilObjectID, cerrs.UnexpectedError(err)
	}
	p.FilmAppearances = fa
	id, err := plRepo.Add(p)
	if err != nil {
		return primitive.NilObjectID, cerrs.UnexpectedError(err)
	}
	return id, nil
}

func (ps PlanetService) getFilmAppearances(planetName string) (qtd uint16, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://swapi.dev/api/planets/?search="+url.QueryEscape(planetName), nil)
	if err != nil {
		return
	}
	req.Close = true
	httpCli := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var data struct {
		Results []struct {
			Films []string `json:"films"`
		} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}
	if len(data.Results) > 0 {
		qtd = uint16(len(data.Results[0].Films))
	}
	return
}

// List lists all planets.
func (ps PlanetService) List() ([]models.Planet, error) {
	plRepo := db.NewPlanet(ps.f.db.Get())
	r, err := plRepo.List()
	if err != nil {
		return nil, cerrs.UnexpectedError(err)
	}
	return r, nil
}

// Get returns a planet by id.
func (ps PlanetService) Get(id string) (models.Planet, error) {
	if err := validate.Exec(validate.NewID(id)); err != nil {
		return models.Planet{}, err
	}
	plRepo := db.NewPlanet(ps.f.db.Get())
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Planet{}, cerrs.UnexpectedError(err)
	}
	r, err := plRepo.FindByID(oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Planet{}, ErrPlanetNotFound
		}
		return models.Planet{}, cerrs.UnexpectedError(err)
	}
	return r, nil
}

// GetByName returns a planet by name.
func (ps PlanetService) GetByName(name string) (models.Planet, error) {
	plRepo := db.NewPlanet(ps.f.db.Get())
	r, err := plRepo.FindByName(strings.TrimSpace(name))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Planet{}, ErrPlanetNotFound
		}
		return models.Planet{}, cerrs.UnexpectedError(err)
	}
	return r, nil
}

// Remove removes a planet by id.
func (ps PlanetService) Remove(id string) error {
	if err := validate.Exec(validate.NewID(id)); err != nil {
		return err
	}
	plRepo := db.NewPlanet(ps.f.db.Get())
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return cerrs.UnexpectedError(err)
	}
	removed, err := plRepo.Remove(oid)
	if err != nil {
		return cerrs.UnexpectedError(err)
	}
	if !removed {
		return ErrPlanetNotFound
	}
	return nil
}
