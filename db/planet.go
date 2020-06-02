package db

import (
	"context"
	"regexp"
	"time"

	"desafio-b2w/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Planet represents a planet repository.
type Planet struct {
	coll *mongo.Collection
}

// NewPlanet creates a planet repository.
func NewPlanet(db *mongo.Database) Planet {
	return Planet{db.Collection("planets")}
}

// Add adds a planet.
func (p Planet) Add(planet models.Planet) (id primitive.ObjectID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	planet.ID = primitive.NewObjectID()
	_, err = p.coll.InsertOne(ctx, planet)
	if err == nil {
		id = planet.ID
	}
	return
}

// List lists all planets.
func (p Planet) List() ([]models.Planet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := p.coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer func() {
		cur.Close(context.TODO())
	}()
	r := make([]models.Planet, 0, 2)
	for cur.Next(context.TODO()) {
		var planet models.Planet
		err = cur.Decode(&planet)
		if err != nil {
			return nil, err
		}
		r = append(r, planet)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return r, nil
}

// FindByID finds a planet by id.
func (p Planet) FindByID(id primitive.ObjectID) (planet models.Planet, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = p.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&planet)
	return
}

// FindByName finds a planet by name.
func (p Planet) FindByName(name string) (planet models.Planet, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = p.coll.FindOne(ctx, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: "^" + regexp.QuoteMeta(name) + "$", Options: "i"}}}).Decode(&planet)
	return
}

// Remove removes a planet by id.
// If the planet is removed, ok is true.
func (p Planet) Remove(id primitive.ObjectID) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := p.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return
	}
	ok = r.DeletedCount > 0
	return
}
