# aws-golang-rds-start-stop-serverless-api

##### AWS Golang serverless API using AWS using SAM to support APIs using API gateway & Lambda to support starting/stopping of RDS instances.
##### support - 100 instances(Postgres & Aurora-Postgres engines)
##### Make sure your CLI session has authorized access to AWS account.

#### 1. Export necessary variables
``` 
    export ORG_ID=foo
    export ENVIRON=uat
    export PROJECT_NAME=play-with-stores
    
```
Note - You do not need to export any variable for local development. You only need to change these variables while deploying it using CICD tools like Codepipeline, Jenkins, TravisCI etc.

#### 2. Deploy locally

```
    make clean build configure run-local
```

#### 3. Deploy on your AWS account

```
    make clean build configure package validate deploy describe outputs
```

You will get the API Gateway URL in the outputs something as below -    
 Ex - https://hahft1bb2c.execute-api.ap-south-1.amazonaws.com/uat

Check your APIs on above URL using curl -

Stop/Start specific instances - 
```
export action=stop #can be stop/start 
export dbidentifier=a #can accept multiple values, comma separated
JSON_STRING='{"type":"'"${action}"'","values":["'"${dbidentifier}"'"]}'

curl -X POST -H 'Content-Type: application/json' \
-H "Authorization: $token" \
--data "$JSON_STRING" \
https://dunhamxtai.execute-api.ap-south-1.amazonaws.com/uat/rdst/action
```

StopAll instances - 
```
JSON_STRING='{"type":"'"${action}"'","values":["'"once"'"]}'

curl -X POST -H 'Content-Type: application/json' \
-H "Authorization: $token" \
--data "$JSON_STRING" \
https://hahft1bb2c.execute-api.ap-south-1.amazonaws.com/uat
```

StartAll instances - 
```
JSON_STRING='{"type":"'"${action}"'","values":["'"once"'"]}'

curl -X POST -H 'Content-Type: application/json' \
-H "Authorization: $token" \
--data "$JSON_STRING" \
https://hahft1bb2c.execute-api.ap-south-1.amazonaws.com/uat

```


Note - You can check application logs in Cloudwatch.


#### 4. Destroy everything

```
    make clean destroy 
```


