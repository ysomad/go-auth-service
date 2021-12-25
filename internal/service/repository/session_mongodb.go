package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/ysomad/go-auth-service/internal/domain"
	apperrors "github.com/ysomad/go-auth-service/pkg/errors"
)

const (
	sessionCollection   = "sessions"
	sessionAccountIDKey = "accountID"
	sessionIDKey        = "_id"
)

type sessionRepo struct {
	*mongo.Collection
}

func NewSessionRepo(db *mongo.Database) *sessionRepo {
	return &sessionRepo{db.Collection(_sessCollection)}
}

func (r *sessionRepo) Create(ctx context.Context, s domain.Session) error {
	ttlIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "createdAt", Value: bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(s.TTL)),
	}

	_, err := r.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		return fmt.Errorf("r.Indexes.CreateOne: %w", err)
	}

	_, err = r.InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("r.InsertOne: %w", err)
	}

	return nil
}

func (r *sessionRepo) Get(ctx context.Context, sid string) (domain.Session, error) {
	var s domain.Session

	if err := r.FindOne(ctx, bson.M{sessionIDKey: sid}).Decode(&s); err != nil {

		if err == mongo.ErrNoDocuments {
			return domain.Session{}, fmt.Errorf("r.FindOne.Decode: %w", apperrors.ErrSessionNotFound)
		}

		return domain.Session{}, fmt.Errorf("r.FindOne.Decode: %w", err)
	}

	return s, nil
}

func (r *sessionRepo) GetAll(ctx context.Context, aid string) ([]domain.Session, error) {
	cursor, err := r.Find(ctx, bson.M{sessionAccountIDKey: bson.M{"$eq": aid}})
	if err != nil {
		return nil, fmt.Errorf("r.Find: %w", err)
	}

	var sessions []domain.Session

	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("cur.All: %w", err)
	}

	return sessions, nil
}

func (r *sessionRepo) Delete(ctx context.Context, sid string) error {
	_, err := r.DeleteOne(ctx, bson.M{sessionIDKey: sid})
	if err != nil {
		return fmt.Errorf("r.DeleteOne: %w", err)
	}

	return nil
}

func (r *sessionRepo) DeleteAll(ctx context.Context, uid string) error {
	_, err := r.DeleteMany(ctx, bson.M{sessionAccountIDKey: uid})
	if err != nil {
		return fmt.Errorf("r.DeleteMany: %w", err)
	}

	return nil
}
