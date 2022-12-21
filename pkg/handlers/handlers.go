package handlers

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bedlad/serverless-stack/pkg/user"
	"github.com/bedlad/serverless-stack/pkg/validators"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	email, exists := req.QueryStringParameters["email"]

	if exists == true {
		if !validators.IsEmailValid(email) {
			log.Panic("Email is invalid")
		}

		result, err := user.DBfetchUser(email, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		} else {
			return apiResponse(http.StatusFound, result)
		}
	} else {
		result, err := user.DBfetchAllUsers(tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		} else {
			return apiResponse(http.StatusFound, result)
		}
	}

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	result, err := user.DBcreateUser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	result, err := user.DBupdateUser(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	email, exists := req.QueryStringParameters["email"]

	if exists == true {
		if !validators.IsEmailValid(email) {
			log.Panic("Email is invalid")
		}

		_, err := user.DBdeleteUser(email, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, nil)
	} else {
		return apiResponse(http.StatusBadGateway, ErrorBody{aws.String("No email provided")})
	}
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, "This method is not allowed")
}
