package user

import (
	"context"
	"database/sql"
	"github.com/SimonMorphy/go-design-pattern/internal/common/config/models"
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type MongoDBUserRepository struct {
	db *mongo.Client
}

func (m MongoDBUserRepository) collection() *mongo.Collection {
	return m.db.Database(viper.GetString("mongo.database")).Collection(viper.GetString("mongo.collection"))
}

type UsrModel struct {
	MongoID   primitive.ObjectID `bson:"_id"`
	ID        uint               `bson:"id" validate:"required"`
	Username  string             `bson:"username" validate:"required,min=3,max=20"`
	Password  string             `bson:"password" validate:"required,min=6,max=32"`
	Email     string             `bson:"email" validate:"required,email"`
	Mobile    string             `bson:"mobile" validate:"omitempty,e164"`
	Address   string             `bson:"address" validate:"omitempty,max=255"`
	Token     string             `bson:"token" validate:"omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt sql.NullTime       `bson:"deleted_at"`
}

func NewMongoDBUserRepository() *MongoDBUserRepository {
	client, err := models.GetClient()
	if err != nil {
		logrus.Panic(err)
	}
	return &MongoDBUserRepository{
		db: client,
	}
}

func (m MongoDBUserRepository) Create(ctx context.Context, user *domain.Usr) (uint, error) {
	model := m.domainToModel(user)
	_, err := m.collection().InsertOne(ctx, model)
	if err != nil {
		logrus.Error(err)
		return -1, err
	}
	return model.ID, nil
}
func (m MongoDBUserRepository) List(ctx context.Context, off, lim int) ([]*domain.Usr, error) {
	ops := &options.FindOptions{}
	ops.SetSkip(int64(off * lim))
	ops.SetLimit(int64(lim))
	cursor, err := m.collection().Find(ctx, bson.M{}, ops)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var userModels []*UsrModel
	if err := cursor.All(ctx, &userModels); err != nil {
		logrus.Error(err)
		return nil, err
	}
	users := make([]*domain.Usr, 0, len(userModels))
	for _, model := range userModels {
		users = append(users, model.ToDomain())
	}
	return users, nil
}
func (m MongoDBUserRepository) Get(ctx context.Context, ID uint) (*domain.Usr, error) {
	_id, _ := primitive.ObjectIDFromHex(strconv.Itoa(int(ID)))
	usr := &UsrModel{}
	err := m.collection().FindOne(ctx, bson.M{"_id": _id}).Decode(usr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if usr == nil {
		return nil, domain.NotFountError{Id: ID}
	}
	return usr.ToDomain(), nil
}

func (m MongoDBUserRepository) Update(ctx context.Context, usr *domain.Usr, fun func(context.Context, *domain.Usr) (*domain.Usr, error)) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBUserRepository) domainToModel(usr *domain.Usr) *UsrModel {
	return &UsrModel{
		MongoID:   primitive.NewObjectID(),
		ID:        usr.ID,
		Username:  usr.Username,
		Password:  usr.Password,
		Email:     usr.Email,
		Mobile:    usr.Mobile,
		Address:   usr.Address,
		Token:     usr.Token,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}
}

func (m UsrModel) ToDomain() *domain.Usr {
	return &domain.Usr{
		Model: gorm.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Username: m.Username,
		Password: m.Password,
		Email:    m.Email,
		Mobile:   m.Mobile,
		Address:  m.Address,
		Token:    m.Token,
	}
}
