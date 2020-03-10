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
```
npm install -g serverless
```

You need to install aws cli 
```
pip install awscli
```

and setup your aws credentials
```
aws configure
```

Of course you need to install [Go](https://golang.org/doc/install)

## Getting started
> By default, a custom authorizer has been enabled for ``create``, ``list``, ``update``, ``delete`` and ``get``. Please replace ``<YOUR_JWT_SECRET_KEY>`` with your JWT Secret Key in serverless.yml (Line 35).  

### Building the code
This script compiles functions to ``bin/handlers/``. 
```
./scripts/build.sh
```

### Deploying to AWS
This script includes the build script and triggers serverless deploy, which will create/update a single CloudFormation stack to provision/update corresponding resources.
```
./scripts/deploy.sh
```

You should see something like 

```bash
************************************************
* Building ...                                  
************************************************
************************************************
* Compiling functions to bin/handlers/ ...      
************************************************
* Compiled authHandler
* Compiled createHandler
* Compiled deleteHandler
* Compiled getHandler
* Compiled listHandler
* Compiled loginHandler
* Compiled updateHandler
************************************************
* Formatting Code ...                           
************************************************
createHandler.go
************************************************
* Build Completed                               
************************************************
************************************************
* Deploying ...                                 
************************************************
Serverless: Packaging service...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Excluding development dependencies...
Serverless: Uploading CloudFormation file to S3...
Serverless: Uploading artifacts...
Serverless: Uploading service auth.zip file to S3 (5.23 MB)...
Serverless: Uploading service list.zip file to S3 (7.2 MB)...
Serverless: Uploading service create.zip file to S3 (7.42 MB)...
Serverless: Uploading service get.zip file to S3 (7.2 MB)...
Serverless: Uploading service login.zip file to S3 (7.23 MB)...
Serverless: Uploading service delete.zip file to S3 (7.14 MB)...
Serverless: Uploading service update.zip file to S3 (7.45 MB)...
Serverless: Validating template...
Serverless: Updating Stack...
Serverless: Checking Stack update progress...
........................
Serverless: Stack update finished...
Service Information
service: serverless-iam-dynamodb
stage: dev
region: ap-southeast-1
stack: serverless-iam-dynamodb-dev
resources: 30
api keys:
  None
endpoints:
  GET - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam
  POST - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam
  PATCH - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam/{id}
  DELETE - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam/{id}
  POST - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam/login
  GET - https://<hash>.execute-api.ap-southeast-1.amazonaws.com/dev/iam/{id}
functions:
  auth: serverless-iam-dynamodb-dev-auth
  list: serverless-iam-dynamodb-dev-list
  create: serverless-iam-dynamodb-dev-create
  update: serverless-iam-dynamodb-dev-update
  delete: serverless-iam-dynamodb-dev-delete
  login: serverless-iam-dynamodb-dev-login
  get: serverless-iam-dynamodb-dev-get
layers:
  None
Serverless: Removing old service artifacts from S3...
Serverless: Run the "serverless" command to setup monitoring, troubleshooting and testing
```