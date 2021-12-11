package dao

import (
	"context"
	"fmt"

	shared_mongo "coolcar/shared/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

type Mongo struct {
	col *mongo.Collection
}

// NewMongo 由 main 函数来传要使用哪个 Database
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("account"),
	}
}

// ResolveAccountID 解析 OpenID，输出 AccountID
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (string, error) {

	// 查询数据库
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, shared_mongo.Set(bson.M{
		openIDField: openID,
	}), options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After))

	// 错误处理
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}

	// 解析出 AccountID，类型为 MongoDB 中的 ObjectID
	var row shared_mongo.ObjID

	// 解码
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}

	// 返回 AccountID
	return row.ID.Hex(), nil

}
