package core

import (
	"context"
	"desafio-b2w/db"
	cerrs "desafio-b2w/errors"
	"desafio-b2w/models"
	"errors"
	"testing"
)

const planetsCollection = "planets"

var planets []models.Planet

func TestPlanets(t *testing.T) {
	f := initForTest()
	f.Run(func(f *Facade) {
		resetCollection(f)
		t.Run("Validations", func(t *testing.T) {
			planetValidations(f, t)
		})
		t.Run("Add", func(t *testing.T) {
			addPlanet(f, &planets[0], t)
		})
		cont := t.Run("Conflict", func(t *testing.T) {
			planetConflict(f, planets[0], t)
		})
		if !cont {
			t.Fatal()
		}
		addPlanet(f, &planets[1], t)
		t.Run("List", func(t *testing.T) {
			listPlanets(f, t)
		})
		t.Run("Get", func(t *testing.T) {
			getPlanet(f, t)
		})
		t.Run("GetByName", func(t *testing.T) {
			getPlanetByName(f, t)
		})
		t.Run("Remove", func(t *testing.T) {
			removePlanet(f, t)
		})
		interrupt(t)
	})
}

func resetCollection(f *Facade) {
	planets = []models.Planet{
		{
			Name:    "Hoth",
			Climate: "frozen",
			Terrain: "tundra, ice caves, mountain ranges",
		},
		{
			Name:    "Dagobah",
			Climate: "murky",
			Terrain: "swamp, jungles",
		},
	}
	f.db.Get().Collection(planetsCollection).Drop(context.TODO())
}

func planetValidations(f *Facade, t *testing.T) {
	_, err := f.Planets().Add(models.Planet{})
	expected := cerrs.InputError("Nome não informado")
	if !errors.Is(err, expected) {
		t.Fatalf("name: expected %#v got %#v", expected, err)
	}

	expected = cerrs.InputError("Clima não informado")
	_, err = f.Planets().Add(models.Planet{Name: "Hoth"})
	if !errors.Is(err, expected) {
		t.Fatalf("climate: expected %#v got %#v", expected, err)
	}

	expected = cerrs.InputError("Terreno não informado")
	_, err = f.Planets().Add(models.Planet{Name: "Hoth", Climate: "frozen"})
	if !errors.Is(err, expected) {
		t.Fatalf("climate: expected %#v got %#v", expected, err)
	}

	resetCollection(f)
}

func addPlanet(f *Facade, pIn *models.Planet, t *testing.T) {
	id, err := f.Planets().Add(*pIn)
	if err != nil {
		t.Fatal(err)
	}
	filmAppearances, err := f.Planets().getFilmAppearances(pIn.Name)
	if err != nil {
		t.Fatal(err)
	}
	pOut, err := db.NewPlanet(f.db.Get()).FindByName(pIn.Name)
	if err != nil {
		t.Fatal(err)
	}
	pIn.ID = id
	if pOut.Name != pIn.Name {
		t.Errorf("planet name want %#v got %#v", pIn.Name, pOut.Name)
	}
	if pOut.Climate != pIn.Climate {
		t.Errorf("planet climate want %#v got %#v", pIn.Climate, pOut.Climate)
	}
	if pOut.Terrain != pIn.Terrain {
		t.Errorf("planet terrain want %#v got %#v", pIn.Terrain, pOut.Terrain)
	}
	if filmAppearances != pOut.FilmAppearances {
		t.Errorf("planet film appearances want %#v got %#v", filmAppearances, pOut.FilmAppearances)
	}
}

func planetConflict(f *Facade, pIn models.Planet, t *testing.T) {
	_, err := f.Planets().Add(pIn)
	if !errors.Is(err, ErrPlanetConflict) {
		t.Fatal("expected planet name conflict")
	}
}

func listPlanets(f *Facade, t *testing.T) {
	list, err := f.Planets().List()
	if err != nil {
		t.Fatal(err)
	}
	gotLen := len(list)
	expectedLen := len(planets)
	if gotLen != 2 {
		t.Errorf("list length: expected %#v got %#v", expectedLen, gotLen)
	}
	for i := 0; i < gotLen; i++ {
		if list[i].Name != planets[i].Name {
			t.Errorf("planet %d: expected %#v got %#v", i, planets[i].Name, list[i].Name)
		}
	}
}

func getPlanet(f *Facade, t *testing.T) {
	for i, p := range planets {
		planet, err := f.Planets().Get(p.ID.Hex())
		if err != nil {
			t.Fatal(err)
		}
		if p.Name != planet.Name {
			t.Errorf("planet %d: expected %#v got %#v", i, p.Name, planet.Name)
		}
	}
}

func getPlanetByName(f *Facade, t *testing.T) {
	planet, err := f.Planets().GetByName("hoth")
	if err != nil {
		t.Fatal(err)
	}
	if planets[0].ID != planet.ID {
		t.Errorf("expected %#v got %#v", planets[0].ID, planet.ID)
	}
}

func removePlanet(f *Facade, t *testing.T) {
	err := f.Planets().Remove(planets[0].ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Planets().Get(planets[0].ID.Hex())
	if !errors.Is(err, ErrPlanetNotFound) {
		t.Errorf("expected error %#v got %#v", ErrPlanetNotFound, err)
	}
}
