# Simple Mock Server #

### Sample Contract
```json
{
  "mappings:": [
    {
      "request": {
        "method": "GET",
        "url": "/hello"
      },
      "responses":  [
           {
             "condition": null,
             "status": 200,
             "fixedDelayMilliseconds": 5000,
             "body": "string content",
             "jsonBody": {
               "status": "Success",
               "message": "Successful response body"
             },
             "headers": {
               "Content-Type": "application/json"
             }
           },
           {
             "condition": {
               "body": {
                 "name": "bar",
                 "surname": "foo"
               }
             },
             "status": 401,
             "fixedDelayMilliseconds": 5000,
             "body": "string content",
             "jsonBody": {
               "status": "Bad Request",
               "message": "Bar foo is not authenticated"
             },
             "headers": {
               "Content-Type": "application/json"
             }
           },
           {
             "condition": {
               "body": {
                 "name": "bar"
               }
             },
             "status": 404,
             "fixedDelayMilliseconds": 5000,
             "body": "string content",
             "jsonBody": {
               "status": "Bad Request",
               "message": "foo not found!"
             },
             "headers": {
               "Content-Type": "application/json"
             }
           }
         ]
    }
  ],
  "auth": {
    "basicAuthCredentials": {
      "username": "jeff@example.com",
      "password": "jeffteenjefftyjeff"
    },
    "tokenCredentials": {
      "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImYzZDVmY2UwLWZiYTMtNDNiOS05NDRjLTMzYmQ1ZTMzNTYwMiIsImlhdCI6MTU4MTM2MTU4OSwiZXhwIjoxNTgxMzY1MTg5fQ.tt11q82zc2i852mEm30YScILqNFP2G_ROnrSZT7Zf28"
    }
  }
}
```
