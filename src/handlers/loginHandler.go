package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-serverless/serverless-iam-dynamodb/src/utils"
)

type Credentials struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type User struct {
	ID            string `json:"id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Age           int    `json:"age,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	Role          string `json:"role,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	ModifiedAt    string `json:"modified_at,omitempty"`
	DeactivatedAt string `json:"deactivated_at,omitempty"`
}

type Response struct {
	Response User `json:"response"`
}

var svc *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	// Initialize a session
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to initialize a session to AWS: %s", err.Error()))
	} else {
		// Create DynamoDB client
		svc = dynamodb.New(session)
	}
}

func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
	)

	// Initialize login
	creds := &Credentials{}

	// Parse request body
	json.Unmarshal([]byte(request.Body), creds)

	result, err := svc.Query(&dynamodb.QueryInput{
		TableName: tableName,
		IndexName: aws.String("IAM_GSI"),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_name": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(creds.UserName),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Got error calling Query:")
		fmt.Println(err.Error())
		// Status Internal Server Error
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	user := User{}

	if len(result.Items) == 0 {
		body, _ := json.Marshal(&Response{
			Response: user,
		})

		// Status OK
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}

	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	match := utils.CheckPasswordHash(creds.Password, user.Password)

	if match {
		body, _ := json.Marshal(&Response{
			Response: user,
		})
		// Status OK
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}

	body, _ := json.Marshal(&Response{
		Response: User{},
	})

	// Status Unauthorized
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 401,
	}, nil

}

func main() {
	lambda.Start(Login)
}
