package dao

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestResovleAccountID(t *testing.T) {
	c := context.Background()

	// 获取客户端对象
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	// 获取 Account 集合管理对象
	m := NewMongo(mc.Database("coolcar"))

	// 解析 OpenID，获取 AccountID
	id, err := m.ResolveAccountID(c, "123")
	if err != nil {
		t.Fatalf("cannot resolve open id: %v", err)
	}

	// 预期值
	want := "61b1e4caf6d536ccefdae779"
	if id != want {
		t.Errorf("resolve account id failed, want: %q, got: %q", want, id)
	}

}
