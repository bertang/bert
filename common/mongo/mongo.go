//@Title mongo.go
//@Description 
//@Author bertang
//@Created bertang 2021/2/17 6:57 下午
package mongo

import (
    "context"
    mongo2 "github.com/bertang/bert/common/config/mongo"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoConn() (*mongo.Client, error) {
    // 设置客户端连接配置
    clientOptions := options.Client().ApplyURI(mongo2.GetMongoConf().String())
    return mongo.Connect(context.TODO(), clientOptions)
}
