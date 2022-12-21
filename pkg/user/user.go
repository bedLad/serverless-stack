package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bedlad/serverless-stack/pkg/validators"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func DBfetchUser(email string, tableName string, dynaClient *dynamodb.DynamoDB) (*User, error) {
	result, err := dynaClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return nil, errors.New("failed to fetch records")
	}

	user := User{}

	if err = dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, errors.New("Failed to Unmarshal Record")
	}

	return &user, nil
}

func DBfetchAllUsers(tableName string, dynaClient *dynamodb.DynamoDB) (*[]User, error) {
	result, err := dynaClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.New("failed to fetch records")
	}

	user := []User{}

	for _, item := range result.Items {
		u := User{}
		if err = dynamodbattribute.UnmarshalMap(item, &u); err != nil {
			return nil, errors.New("Failed to Unmarshal Record")
		}

		user = append(user, u)
	}

	return &user, nil
}

func DBcreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*User, error) {
	user := User{}
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New("Error while unmarshalling the data")
	}

	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New("Invalid email provided")
	}

	currUser, _ := DBfetchUser(user.Email, tableName, dynaClient)
	if currUser != nil && len(currUser.Email) != 0 {
		return nil, errors.New("User already exists")
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New("Error while marshalling the data")
	}

	_, err = dynaClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.New("Error while inserting the data")
	}

	return &user, err
}

func DBupdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*User, error) {
	user := User{}
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New("Error while unmarshalling the data")
	}

	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New("Invalid email provided")
	}

	currUser, _ := DBfetchUser(user.Email, tableName, dynaClient)
	if currUser == nil {
		return nil, errors.New("User does not exist")
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New("Error while marshalling the data")
	}

	_, err = dynaClient.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.New("Error while inserting the data")
	}

	return &user, err
}

func DBdeleteUser(email string, tableName string, dynaClient *dynamodb.DynamoDB) (string, error) {
	currUser, _ := DBfetchUser(email, tableName, dynaClient)
	if currUser == nil {
		return "nil", errors.New("User does not exist")
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return "nil", errors.New("Error while deleting the data")
	}

	return "User successfully deleted", err
}
