// Copyright 2018 Kuei-chun Chen. All rights reserved.

package examples

import (
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/simagix/keyhole/mdb"
)

/*
 * https://docs.mongodb.com/manual/reference/operator/aggregation/graphLookup/
 */
func TestAggregateGraphLookup(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M

	client = getMongoClient()
	seedCarsData(client, dbName)

	pipeline := `
	[{
		"$graphLookup": {
			"from": "employees",
			"startWith": "$manager",
			"connectFromField": "manager",
			"connectToField": "_id",
			"as": "employeeHierarchy"
		}
	}]
	`

	collection = client.Database(dbName).Collection("employees")
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, mdb.MongoPipeline(pipeline), opts); err != nil {
		t.Fatal(err)
	}
	defer cur.Close(ctx)
	count := int64(0)
	for cur.Next(ctx) {
		cur.Decode(&doc)
		count++
	}

	if 0 == count {
		t.Fatal("no doc found")
	}
}
