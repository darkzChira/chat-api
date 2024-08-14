package repository

import (
	"chat-api/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Print("FindUserByUsername")
			return nil, err
		}

		log.Print(err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"online_status": user.OnlineStatus,
		},
	}
	_, err = r.collection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *UserRepository) FindAllExceptUserID(ctx context.Context, userID string) ([]models.User, error) {
	objectId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var users []models.User
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$ne": objectId}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
