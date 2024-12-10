package repositories

import (
	"context"
	"crawl-manager-backend/app/models"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepositoryImpl struct {
	DB *mongo.Client
}

var userModel models.User

func NewAuthRepository(db *mongo.Client) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{DB: db}
}

// Create a unique index on email during initialization
func (r *AuthRepositoryImpl) CreateEmailUniqueIndex() error {
	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a unique index on the email field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *AuthRepositoryImpl) FindUserByUsername(username string) (*models.User, error) {
	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) CreateUser(user *models.User) error {
	if r.DB == nil {
		return errors.New("DB is nil")
	}

	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email already exists
	existingUser, err := r.FindUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Check if username already exists
	existingUsername, err := r.FindUserByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUsername != nil {
		return errors.New("username already exists")
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		// Check for duplicate key error (duplicate email)
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email or username is already in use")
		}
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

func (r *AuthRepositoryImpl) UpdateUser(username string, updatedUser *models.User) error {
	if r.DB == nil {
		return errors.New("DB is nil")
	}

	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	update := bson.M{"$set": updatedUser}

	// Check for potential email conflicts during update
	if updatedUser.Email != "" {
		existingUser, err := r.FindUserByEmail(updatedUser.Email)
		if err != nil {
			return err
		}
		if existingUser != nil && existingUser.Username != username {
			return errors.New("email already in use by another user")
		}
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no user found with the given username")
	}

	return nil
}

func (r *AuthRepositoryImpl) FindUserByID(userID string) (*models.User, error) {
	collection := r.DB.Database(DBName).Collection(userModel.GetTableName())

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
