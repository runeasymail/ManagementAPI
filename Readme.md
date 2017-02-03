# API for runeasymail.com  
[![Build Status](https://travis-ci.org/runeasymail/ManagementAPI.svg?branch=master)](https://travis-ci.org/runeasymail/ManagementAPI)

## Endpoints

### POST /auth
For all of methods above you need to send `Auth-token` header with token which can be obtain from call to `/auth
 method
 
Example:
```curl
curl -X POST  -F "username=yuksel" -F "password=test" "http://localhost:8081/auth"
```
If auth is successfull result must be like:
```json
{"result":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODYxMTUzNTEsInVzZXJuYW1lIjoieXVrc2VsIn0.jIaigJSR8NQOuXR6Sjrcccz_q9KrodD2rQ6e0dLWlgE"}
```

### GET /domains
It's will give list of all availible domains with ther internal id

Example of `curl` request:
```bash
curl -X GET -H "Auth-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODYxMTQ2MzMsInVzZXJuYW1lIjoieXVrc2VsIn0.AaQxjRU7PLT9A-CMyYjcXWEki3oQxA8GUv8rANEj59M" -H "Cache-Control: no-cache" "http://localhost:8081/domains"
```

Example result:
```json
{
   "domains":[
      {
         "id":13,
         "name":"yuks212.me",
         "users_count":1
      },
      {
         "id":12,
         "name":"yuks21.me",
         "users_count":1
      },
      {
         "id":11,
         "name":"yuks2.me",
         "users_count":1
      },
      {
         "id":10,
         "name":"yuks.me",
         "users_count":1
      },
      {
         "id":1,
         "name":"mail.yuks.me",
         "users_count":1
      }
   ]
}
```

### POST /domains
params:
```
domain - string
username - string (email)
password - string
```
It will create a new domain and it will add new user to it

### GET /users/:domain_id
Will return list with users on specific domain id

Example:
```json
{
   "users":[
      {
         "id":1,
         "domain_id":1,
         "email":"admin@mail.yuks.me"
      }
   ]
}
```

### POST /users/:domain_id
params:
```
domain_id - id of domain
email - email address of new user
password - password of new user (clear) it will encrypted
```

## Instalation
```bash
cd ...
go build -o main
```

## 

### Used libraries
* github.com/gin-gonic/gin
* github.com/go-sql-driver/mysql
* github.com/go-ini/ini 
* github.com/op/go-logging
* github.com/jmoiron/sqlx
* github.com/dgrijalva/jwt-go

