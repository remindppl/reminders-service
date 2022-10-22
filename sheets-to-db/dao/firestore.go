package dao

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
)

var (
	envProjectId = "ENV_PROJECT_ID"
)

type Dao interface {
	Put(ctx context.Context, req *FollowUpRequest) error
	Get(ctx context.Context, id string) (*FollowUpRequest, error)
	Close(ctx context.Context) error
}

type FirestoreDao struct {
	client     *firestore.Client
	collection string
}

func (db *FirestoreDao) Put(ctx context.Context, req *FollowUpRequest) error {
	docRef := db.client.Doc(db.collection + "/" + req.RequestID)
	_, err := docRef.Set(ctx, req)
	return err
}

func (db *FirestoreDao) Get(ctx context.Context, id string) (*FollowUpRequest, error) {
	docRef := db.client.Doc(db.collection + "/" + id)
	ds, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	out := FollowUpRequest{}
	if err := ds.DataTo(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (db *FirestoreDao) Close(ctx context.Context) error {
	return db.client.Close()
}

func New(ctx context.Context, collection string) (*FirestoreDao, error) {
	projectID, found := os.LookupEnv(envProjectId)
	if !found {
		panic("Didn't find projectID env var: " + envProjectId)
	}
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client. Err: %v", err)
	}
	return &FirestoreDao{
		client:     c,
		collection: collection,
	}, nil
}
