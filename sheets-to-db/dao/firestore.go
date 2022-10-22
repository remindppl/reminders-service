package dao

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var (
	envProjectId = "ENV_PROJECT_ID"
)

type Dao interface {
	Put(ctx context.Context, req *FollowUpRequest) error
	Get(ctx context.Context, id string) (*FollowUpRequest, error)
}

type FirestoreDao struct {
	client     *firestore.Client
	collection string
}

func (db *FirestoreDao) Put(ctx context.Context, req *FollowUpRequest) error {
	docRef := db.client.Collection(db.collection).NewDoc()
	_, err := docRef.Create(ctx, req)
	return err
}

func (db *FirestoreDao) Get(ctx context.Context, id string) (*FollowUpRequest, error) {
	docRef := db.client.Collection(db.collection).Doc(id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	var followup FollowUpRequest
	doc.DataTo(&followup)
	return &followup, nil
}

func New(ctx context.Context, collection string) *FirestoreDao {
	projectID, found := os.LookupEnv(envProjectId)
	if !found {
		panic("Didn't find projectID env var: " + envProjectId)
	}

	c, _ := firestore.NewClient(ctx, projectID)
	return &FirestoreDao{
		client:     c,
		collection: collection,
	}
}
