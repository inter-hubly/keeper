package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/app/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact interface {
	SaveContact(ctx context.Context, contact *domain.Contact) (string, error)
	FindContacts(ctx context.Context, tenant string) ([]domain.Contact, error)
}

type contactRepository struct {
	collection string
	connection hmongo.NoSqlConn
}

var (
	contactRepositoryOnce sync.Once
	contact               *contactRepository
)

func NewContact(ctx context.Context) *contactRepository {
	contactRepositoryOnce.Do(func() {
		contact = &contactRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "contact",
		}
	})
	return contact
}

func (r *contactRepository) SaveContact(ctx context.Context, contact *domain.Contact) (string, error) {
	one, err := r.connection.GetCollection(ctx, r.collection).InsertOne(ctx, contact)
	if err != nil {
		hlog.Error(ctx, "contactRepository.SaveContact", err.Error())
		return "", err
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *contactRepository) FindContacts(ctx context.Context, tenant string) ([]domain.Contact, error) {
	hlog.Debug(ctx, "contactRepository.FindContacts", "find contacts")

	cur, err := r.connection.GetCollection(ctx, r.collection).Find(ctx, bson.M{
		"tenantId": tenant,
	})

	if err != nil {
		hlog.Error(ctx, "contactRepository.FindContacts", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	contacts := make([]domain.Contact, 0)
	for cur.Next(ctx) {
		var contactDb domain.Contact
		if err = cur.Decode(&contactDb); err != nil {
			hlog.Error(ctx, "contactRepository.FindContacts", err.Error())
			return nil, err
		}
		contacts = append(contacts, contactDb)
	}
	return contacts, nil
}
