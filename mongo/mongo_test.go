package mongo

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestWriteErrorsOnDisconnectedClient(t *testing.T) {
	// NewClient does not connect. Each operation fails locally before any I/O.
	client, err := driver.NewClient(options.Client().ApplyURI("mongodb://example.invalid:27017"))
	if err != nil {
		t.Fatal(err)
	}
	collection := &mongoCollection{coll: client.Database("synthetic").Collection("items")}
	ctx := context.Background()

	t.Run("InsertOne", func(t *testing.T) {
		id, err := collection.InsertOne(ctx, bson.M{"title": "synthetic"})
		if !errors.Is(err, driver.ErrClientDisconnected) {
			t.Fatalf("error = %v, want ErrClientDisconnected", err)
		}
		if id != nil {
			t.Fatalf("id = %v, want nil", id)
		}
	})
	t.Run("InsertMany", func(t *testing.T) {
		ids, err := collection.InsertMany(ctx, []interface{}{bson.M{"title": "synthetic"}})
		if !errors.Is(err, driver.ErrClientDisconnected) {
			t.Fatalf("error = %v, want ErrClientDisconnected", err)
		}
		if ids != nil {
			t.Fatalf("ids = %v, want nil", ids)
		}
	})
	t.Run("DeleteOne", func(t *testing.T) {
		count, err := collection.DeleteOne(ctx, bson.M{"title": "synthetic"})
		if !errors.Is(err, driver.ErrClientDisconnected) {
			t.Fatalf("error = %v, want ErrClientDisconnected", err)
		}
		if count != 0 {
			t.Fatalf("count = %d, want 0", count)
		}
	})
}
