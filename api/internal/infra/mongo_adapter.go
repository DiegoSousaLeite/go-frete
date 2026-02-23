package infra

import (
	"context"
	"time"

	"go-frete/api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
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

// SaveHistory implementa a interface domain.ConversionSaver
func (m *MongoDBAdapter) SaveHistory(record domain.ConversionRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insere a struct que será traduzida para BSON (formato do Mongo)
	_, err := m.database.Collection(conversionHistory).InsertOne(ctx, record)
	return err
}

func (m *MongoDBAdapter) GetLastConversions(limit int) ([]domain.ConversionRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Configura a ordenação (data em ordem decrescente: -1) e aplica o limite
	opts := options.Find().
		SetSort(bson.D{{Key: "data", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := m.database.Collection(conversionHistory).Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.ConversionRecord
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetConversionsByCurrency busca as conversões de uma moeda ordenadas da mais antiga para a mais nova
func (m *MongoDBAdapter) GetConversionsByCurrency(currency string) ([]domain.ConversionRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Filtro: Onde a chave "currency" seja igual a moeda pedida
	filter := bson.D{{Key: "currency", Value: currency}}

	// Ordenação: data crescente (1)
	opts := options.Find().SetSort(bson.D{{Key: "data", Value: 1}})

	cursor, err := m.database.Collection(conversionHistory).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.ConversionRecord
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
