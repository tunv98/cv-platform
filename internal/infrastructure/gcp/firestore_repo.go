package gcp

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cv-platform/internal/domain"
	"cv-platform/internal/port"
)

type FirestoreCVRepo struct {
	cl   *firestore.Client
	coll string
}

func NewFirestoreCVRepo(ctx context.Context, projectID string, credsJSON []byte) (port.CVRepository, error) {
	var (
		cl  *firestore.Client
		err error
	)
	if len(credsJSON) > 0 {
		cl, err = firestore.NewClient(ctx, projectID, option.WithCredentialsJSON(credsJSON))
	} else {
		cl, err = firestore.NewClient(ctx, projectID)
	}
	if err != nil {
		return nil, err
	}
	return &FirestoreCVRepo{cl: cl, coll: "cvs"}, nil
}

func (r *FirestoreCVRepo) Create(cv *domain.CV) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.cl.Collection(r.coll).Doc(cv.ID).Create(ctx, cv)
	return err
}

func (r *FirestoreCVRepo) Update(cv *domain.CV) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.cl.Collection(r.coll).Doc(cv.ID).Set(ctx, cv)
	return err
}

func (r *FirestoreCVRepo) FindByID(id string) (*domain.CV, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc, err := r.cl.Collection(r.coll).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var cv domain.CV
	if err := doc.DataTo(&cv); err != nil {
		return nil, err
	}
	return &cv, nil
}

func (r *FirestoreCVRepo) List(limit int, cursor string) ([]domain.CV, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.cl.Collection(r.coll).OrderBy("CreatedAt", firestore.Desc).Limit(limit)
	it := q.Documents(ctx)
	var out []domain.CV
	for {
		doc, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, "", err
		}
		var cv domain.CV
		if err := doc.DataTo(&cv); err != nil {
			return nil, "", err
		}
		out = append(out, cv)
	}
	// Simplified: no real cursor
	return out, "", nil
}
