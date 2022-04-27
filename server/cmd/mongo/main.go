package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	client, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017"))
	if err != nil {
		panic(err)
	}
	col := client.Database("coolcar").Collection("account")
	//insertRows(c, col)
	findRows(c, col)
}
func insertRows(c context.Context, col *mongo.Collection) {
	many, err := col.InsertMany(c, []interface{}{
		bson.M{
			"open_id": "123",
		},
		bson.M{
			"open_id": "321",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(many)

}
func findRows(c context.Context, col *mongo.Collection) {
	res := col.FindOne(c, bson.M{
		"open_id": "123",
	})
	fmt.Printf("%+v\n", res)
	var row struct {
		ID     primitive.ObjectID `bson:"_id"`
		OpenID string             `bson:"open_id"`
	}
	err := res.Decode(&row)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", row) //%+v先输出字段名字，再输出该字段的值 ,%#v 先输出结构体名字值，再输出结构体（字段名字+字段的值）

	find, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for find.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err := find.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}
