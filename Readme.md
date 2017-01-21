# API for runeasymail.com  
[![Build Status](https://travis-ci.org/runeasymail/ManagementAPI.svg?branch=master)](https://travis-ci.org/runeasymail/ManagementAPI)

## Endpoints

### GET /domains
It's will give list of all availible domains with ther internal id

Example:
```
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
```
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
```
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

