package infra

import (
	"context"
	"time"

	"go-frete/api/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const conversionHistory = "conversion_history"

type MongoDBAdapter struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBAdapter conecta no banco e retorna o adapter
func NewMongoDBAdapter(uri, dbName string) (*MongoDBAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDBAdapter{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// SaveHistory implementa a interface domain.ConversionRepository
func (m *MongoDBAdapter) SaveHistory(record domain.ConversionRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insere a struct que será traduzida para BSON (formato do Mongo)
	_, err := m.database.Collection(conversionHistory).InsertOne(ctx, record)
	return err
}
