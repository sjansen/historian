package storage

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/oklog/ulid"

	"github.com/sjansen/historian/internal/dto"
)

var twoWeeks = time.Hour * 24 * 14
var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

type DynamoDBRepo struct {
	svc   *dynamodb.DynamoDB
	table *string
}

func NewDynamoDBRepo(table string) (*DynamoDBRepo, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	repo := &DynamoDBRepo{
		svc:   dynamodb.New(sess),
		table: aws.String(table),
	}
	return repo, err
}

func (repo *DynamoDBRepo) Add(msg *dto.Message) (id string, err error) {
	now := time.Now().UTC()
	id = ulid.MustNew(ulid.Timestamp(now), entropy).String()
	expires := strconv.FormatInt(now.Add(twoWeeks).Unix(), 10)
	timestamp := strconv.FormatInt(msg.Timestamp.Unix(), 10)
	input := &dynamodb.PutItemInput{
		TableName: repo.table,
		Item: map[string]*dynamodb.AttributeValue{
			"event-id":  {S: aws.String(id)},
			"timestamp": {N: aws.String(timestamp)},
			"expires":   {N: aws.String(expires)},
			"data":      {S: aws.String(msg.RawData)},
		},
	}
	if _, err := repo.svc.PutItem(input); err != nil {
		return "", err
	}
	return id, nil
}
