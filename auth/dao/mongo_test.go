package dao

import (
	"context"
	"os"
	"testing"

	shared_mongo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string

func TestResovleAccountID(t *testing.T) {
	c := context.Background()

	// 获取客户端对象
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	// 获取 Account 集合管理对象
	m := NewMongo(mc.Database("coolcar"))

	// 插入表格数据，做表格测试
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			shared_mongo.IDField: MustObjID("61b1e4caf6d536ccefdae779"),
			openIDField:          "open_id_123",
		},
		bson.M{
			shared_mongo.IDField: MustObjID("61b1e4caf6d536ccefdae778"),
			openIDField:          "open_id_456",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}

	// 不存在的用户要插入的 ObjectID
	m.newObjIDFunc = func() primitive.ObjectID {
		return MustObjID("61b1e4caf6d536ccefdae777")
	}

	// 测试样例
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "open_id_123",
			want:   "61b1e4caf6d536ccefdae779",
		},
		{
			name:   "another_existing_user",
			openID: "open_id_456",
			want:   "61b1e4caf6d536ccefdae778",
		},
		{
			name:   "new_user",
			openID: "open_id_789",
			want:   "61b1e4caf6d536ccefdae777",
		},
	}

	// 测试
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), testCase.openID)
			if err != nil {
				t.Errorf("cannot resolve id for %q, error: %v", testCase.openID, err)
			}
			if id != testCase.want {
				t.Errorf("Error!!! resolve account id failed, want: %q, got %q", testCase.want, id)
			}
		})
	}
}

// 必须这么命名
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}

// 生成固定的 ObjectID
func MustObjID(hex string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(hex)
	return objID
}
