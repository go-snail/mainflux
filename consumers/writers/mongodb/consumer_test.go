// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mongodb_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mainflux/mainflux/consumers/writers/mongodb"
	"github.com/mainflux/mainflux/pkg/transformers/senml"

	log "github.com/mainflux/mainflux/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	port        string
	addr        string
	testLog, _  = log.New(os.Stdout, log.Info.String())
	testDB      = "test"
	collection  = "messages"
	db          mongo.Database
	msgsNum     = 100
	valueFields = 5
	subtopic    = "topic"
)

var (
	v       float64 = 5
	stringV         = "value"
	boolV           = true
	dataV           = "base64"
	sum     float64 = 42
)

func TestSave(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(addr))
	require.Nil(t, err, fmt.Sprintf("Creating new MongoDB client expected to succeed: %s.\n", err))

	db := client.Database(testDB)
	repo := mongodb.New(db)

	now := time.Now().Unix()
	msg := senml.Message{
		Channel:    "45",
		Publisher:  "2580",
		Protocol:   "http",
		Name:       "test name",
		Unit:       "km",
		Time:       13451312,
		UpdateTime: 5456565466,
	}
	var msgs []senml.Message

	for i := 0; i < msgsNum; i++ {
		// Mix possible values as well as value sum.
		count := i % valueFields
		switch count {
		case 0:
			msg.Subtopic = subtopic
			msg.Value = &v
		case 1:
			msg.BoolValue = &boolV
		case 2:
			msg.StringValue = &stringV
		case 3:
			msg.DataValue = &dataV
		case 4:
			msg.Sum = &sum
		}

		msg.Time = float64(now + int64(i))
		msgs = append(msgs, msg)
	}

	err = repo.Consume(msgs)
	assert.Nil(t, err, fmt.Sprintf("Save operation expected to succeed: %s.\n", err))

	count, err := db.Collection(collection).CountDocuments(context.Background(), bson.D{})
	assert.Nil(t, err, fmt.Sprintf("Querying database expected to succeed: %s.\n", err))
	assert.Equal(t, int64(msgsNum), count, fmt.Sprintf("Expected to have %d value, found %d instead.\n", msgsNum, count))
}
