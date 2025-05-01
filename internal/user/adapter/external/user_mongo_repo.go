package external

import (
	"context"
	"database/sql"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

type MongoDBUserRepository struct {
	db *mongo.Client
}

func (m MongoDBUserRepository) Delete(ctx context.Context, ID uint) error {
	one, err := m.collection().DeleteOne(ctx, bson.M{"id": ID})
	if err != nil && one.DeletedCount == 0 {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return err
	}
	return nil
}

func (m MongoDBUserRepository) collection() *mongo.Collection {
	return m.db.Database(viper.GetString("mongo.database")).Collection(viper.GetString("mongo.collection"))
}

type UsrModel struct {
	MongoID   primitive.ObjectID `bson:"_id,omitempty"`
	ID        uint               `bson:"id" validate:"omitempty"`
	Username  string             `bson:"username,omitempty" validate:"omitempty,min=3,max=20"`
	Password  string             `bson:"password,omitempty" validate:"omitempty,min=6,max=32"`
	Email     string             `bson:"email,omitempty" validate:"omitempty,email"`
	Mobile    string             `bson:"mobile,omitempty" validate:"omitempty,e164"`
	Address   string             `bson:"address,omitempty" validate:"omitempty,max=255"`
	Token     string             `bson:"token,omitempty" validate:"omitempty"`
	CreatedAt *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt sql.NullTime       `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func NewMongoDBUserRepository() (*MongoDBUserRepository, func()) {
	client, cleanUp, err := models.InitMongoDB()
	if err != nil {
		logrus.Panic(err)
	}
	return &MongoDBUserRepository{
			db: client,
		}, func() {
			cleanUp("mongo")
		}
}

func (m MongoDBUserRepository) Create(ctx context.Context, user *domain.Usr) (uint, error) {
	user.ID = uint(uuid.New().ID())
	model := m.domainToModel(user)
	model.MongoID = primitive.NewObjectID()
	_, err := m.collection().InsertOne(ctx, model)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return model.ID, nil
}
func (m MongoDBUserRepository) List(ctx context.Context, off, lim int) ([]*domain.Usr, error) {
	ops := &options.FindOptions{}
	ops.SetSkip(int64((off - 1) * lim))
	ops.SetLimit(int64(lim))
	cursor, err := m.collection().Find(ctx, bson.M{}, ops)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var userModels []*UsrModel
	if err = cursor.All(ctx, &userModels); err != nil {
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

	usr := &UsrModel{}
	err := m.collection().FindOne(ctx, bson.M{"id": ID}).Decode(usr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if usr == nil {
		return nil, errors.NewWithError(errors.ErrnoUserNotFoundError, domain.NotFountError{Id: ID})
	}
	return usr.ToDomain(), nil
}

func (m MongoDBUserRepository) Update(ctx context.Context, usr *domain.Usr, fun func(context.Context, *domain.Usr) (*domain.Usr, error)) error {
	_usr, err := fun(ctx, usr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	model := m.domainToModel(_usr)
	_, err = m.collection().UpdateOne(ctx, bson.M{"id": model.ID}, bson.D{{"$set", model}})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil

}

func (m MongoDBUserRepository) domainToModel(usr *domain.Usr) *UsrModel {
	return &UsrModel{
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
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Password:  m.Password,
		Email:     m.Email,
		Mobile:    m.Mobile,
		Address:   m.Address,
		Token:     m.Token,
		DeletedAt: gorm.DeletedAt(m.DeletedAt),
	}
}
