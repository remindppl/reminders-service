package dao

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

var (
	envProjectId = "env-project-id"
)

type Dao interface {
	Put(ctx context.Context, req *FollowUpRequest) error
	Get(ctx context.Context, id string) (*FollowUpRequest, error)
}

type FirestoreDao struct {
	client *firestore.Client
}

func (db *FirestoreDao) Put(ctx context.Context, req *FollowUpRequest) error {
	docRef := db.client.Collection("collection").NewDoc()
	docRef.ID = req.RequestID
	_, err := docRef.Create(ctx, req)
	return err
}

func (db *FirestoreDao) Get(ctx context.Context, id string) (*FollowUpRequest, error) {
	docRef := db.client.Collection("collection").Doc(id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	var followup FollowUpRequest
	doc.DataTo(&followup)
	return &followup, nil
}

func New(ctx context.Context) FirestoreDao {
	projectID, found := os.LookupEnv(envProjectId)
	if !found {
		panic("Didn't find projectID env var")
	}

	c, _ := firestore.NewClient(ctx, projectID)
	return FirestoreDao{
		client: c,
	}
}
