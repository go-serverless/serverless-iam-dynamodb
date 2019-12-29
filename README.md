# serverless-iam-dynamodb
Building serverless CRUD services in Go with DynamoDB

## Project structure
/.serverless 

It will be created automatically when running ``serverless deploy`` in where deployment zip files, cloudformation stack files will be generated

/bin

This is the folder where our built Go codes are placed

/scripts

General scripts for building Go codes and deployment

/src/handlers

All Lambda handlers will be placed here

/src/utils

General functions go here

## Prerequisites

You need to install serverless cli
````
npm install -g serverless
````

You need to install aws cli 
````
pip install awscli
````

and setup your aws credentials
````
aws configure
````

Of course you need to install [Go](https://golang.org/doc/install)

## Building the code
````
./scripts/build.sh
````

## Deploying to AWS
````
./scripts/deploy.sh
````