package orm

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database interface {
	FindOne(c context.Context, model interface{}, filter interface{}, item interface{}) error
	FindMany(c context.Context, model interface{}, filter interface{}, items interface{}) error
	FindPage(c context.Context, model interface{}, filter interface{}, page, size int, items interface{}) error
	InsertOne(c context.Context, model interface{}, item interface{}) error
	DeleteOne(c context.Context, model interface{}, item interface{}) error
	UpdateOne(c context.Context, model interface{}, filter interface{}, update interface{}) error
	UpsertOne(c context.Context, model interface{}, update map[string]interface{}, create interface{}) error
	Transaction(c context.Context, fn func(tx *gorm.DB) error) error
	Raw(c context.Context) *gorm.DB
}

type database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) Database {
	return &database{db: db}
}

func (dao *database) Raw(c context.Context) *gorm.DB {
	return dao.db
}

func (dao *database) FindOne(c context.Context, model interface{}, filter interface{}, item interface{}) error {
	return dao.db.WithContext(c).Model(model).Where(filter).First(item).Error
}

func (dao *database) FindMany(c context.Context, model interface{}, filter interface{}, items interface{}) error {
	return dao.db.WithContext(c).Model(model).Where(filter).Find(items).Error
}

func (dao *database) FindPage(c context.Context, model interface{}, filter interface{}, page, size int, items interface{}) error {
	return dao.db.WithContext(c).Model(model).Where(filter).Offset((page - 1) * size).Limit(size).Find(items).Error
}

func (dao *database) InsertOne(c context.Context, model interface{}, item interface{}) error {
	return dao.db.WithContext(c).Model(model).Create(item).Error
}

func (dao *database) DeleteOne(c context.Context, model interface{}, item interface{}) error {
	return dao.db.WithContext(c).Model(model).Where(item).Delete(item).Error
}

func (dao *database) UpdateOne(c context.Context, model interface{}, filter interface{}, update interface{}) error {
	return dao.db.WithContext(c).Model(model).Where(filter).Updates(update).Error

}

func (dao *database) UpsertOne(c context.Context, model interface{}, update map[string]interface{}, create interface{}) error {
	err := dao.db.WithContext(c).Model(model).Clauses(clause.OnConflict{
		//Columns:  // mysql 可以不写
		DoUpdates: clause.Assignments(update),
	}).Create(create).Error
	return err
}

func (dao *database) Transaction(c context.Context, fn func(tx *gorm.DB) error) error {
	return dao.db.WithContext(c).Transaction(fn)
}
