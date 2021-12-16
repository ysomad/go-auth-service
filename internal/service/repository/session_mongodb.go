package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/ysomad/go-auth-service/internal/entity"
)

const sessionCollection = "sessions"

type sessionRepo struct {
	*mongo.Collection
}

func NewSessionRepo(db *mongo.Database) *sessionRepo {
	return &sessionRepo{db.Collection(sessionCollection)}
}

// Create creates new user session in redis
func (r *sessionRepo) Create(ctx context.Context, s entity.Session) error {
	ttlIndex := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "createdAt", Value: bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(s.TTL.Seconds())),
	}

	_, err := r.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		return fmt.Errorf("r.Indexes.CreateOne: %w", err)
	}

	primitive.NewObjectID()
	_, err = r.InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("r.InsertOne: %w", err)
	}

	return nil
}

func (r *sessionRepo) Get(ctx context.Context, sid string) (entity.Session, error) {
	var s entity.Session

	if err := r.FindOne(ctx, bson.M{"_id": sid}).Decode(&s); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Session{}, fmt.Errorf("r.FindOne.Decode: %w", entity.ErrSessionNotFound)
		}

		return entity.Session{}, fmt.Errorf("r.FindOne.Decode: %w", err)
	}

	return s, nil
}

func (r *sessionRepo) GetAll(ctx context.Context, uid string) ([]entity.Session, error) {
	cur, err := r.Find(ctx, bson.M{"userID": bson.M{"$eq": uid}})
	if err != nil {
		return nil, fmt.Errorf("r.Find: %w", err)
	}

	var sessions []entity.Session

	if err = cur.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("cur.All: %w", err)
	}

	return sessions, nil
}

func (r *sessionRepo) Delete(ctx context.Context, sid string) error {
	_, err := r.DeleteOne(ctx, bson.M{"_id": sid})
	if err != nil {
		return fmt.Errorf("r.DeleteOne: %w", err)
	}

	return nil
}

func (r *sessionRepo) DeleteAll(ctx context.Context, uid string) error {
	_, err := r.DeleteMany(ctx, bson.M{"userID": uid})
	if err != nil {
		return fmt.Errorf("r.DeleteMany: %w", err)
	}

	return nil
}
