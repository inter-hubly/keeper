//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

package contact

import (
	"context"
	"sync"

	"github.com/inter-hubly/keeper/internal/domain"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	SaveContact(ctx context.Context, contact *domain.Contact) (string, error)
	FindContacts(ctx context.Context, tenant string) ([]domain.Contact, error)
}

type contactRepository struct {
	collection string
	connection hmongo.NoSqlConn
}

var (
	_contactRepositoryOnce sync.Once
	_contactRepository     *contactRepository
)

func NewContact(ctx context.Context) *contactRepository {
	_contactRepositoryOnce.Do(func() {
		_contactRepository = &contactRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "contact",
		}
	})
	return _contactRepository
}

func (r *contactRepository) SaveContact(ctx context.Context, contact *domain.Contact) (string, error) {
	one, err := r.connection.GetCollection(ctx, r.collection).InsertOne(ctx, contact)
	if err != nil {
		hlog.Error(ctx, "contact.repository.SaveContact", err.Error())
		return "", err
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *contactRepository) FindContacts(ctx context.Context, tenant string) ([]domain.Contact, error) {
	hlog.Debug(ctx, "contact.repository.FindContacts", "find contacts")

	cur, err := r.connection.GetCollection(ctx, r.collection).Find(ctx, bson.M{
		"tenantId": tenant,
	})

	if err != nil {
		hlog.Error(ctx, "contact.repository.FindContacts", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	contacts := make([]domain.Contact, 0)
	for cur.Next(ctx) {
		var contactDb domain.Contact
		if err = cur.Decode(&contactDb); err != nil {
			hlog.Error(ctx, "contact.repository.FindContacts", err.Error())
			return nil, err
		}
		contacts = append(contacts, contactDb)
	}
	return contacts, nil
}
