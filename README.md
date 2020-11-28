# MyMoment Backend

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. This app uses go version go version `go1.15 windows/amd64`

### Prerequisites
##### Go [[link](https://golang.org/dl)]
Version I am using : `go1.15 windows/amd64`

##### Redis [[link](https://redis.io/download)]
Version I am using : `Redis server version 2.4.5`

##### MySQL DB [[link](https://www.mysql.com/downloads)]
Version I am using : `mysql  Ver 8.0.16 for Win64 on x86_64 (MySQL Community Server - GPL)`

##### MailHog [[link](https://github.com/mailhog/MailHog)]

## Project setup

### Add Config File
Create new file `config.json` inside `config` folder.

The following is brief axplanation of the config
```
{
    "DB":{
        "Host": "<your db host, example : localhost>",
        "Port": <your db port, default port for mysql : 3306>,
        "Username": "<your db username>",
        "Password":"<your db password>",
        "DbName":"<your db name>"
    },
    "SMTP":{
        "From":"<smtp from address, you can type any email address for development environment>",
        "Host":"<smtp host, example : localhost >",
        "Port":<smtp port, MailHog default port : 1025>,
        "Username":"<smtp username, can be left empty for development>",
        "Password":"<smtp password, can be left empty for development>",
    },
    "EmailVerification":{
        "Subject":"<subject to be used for email verification>"
    },
    "ResetPassword":{
        "Subject":"<subject to be used for reset password email>"
    },
    "Redis":{
        "Host":"<redis host, example : localhost>",
        "Password":"<redis passwrod, can be left empty for development>",
        "Port":<redis port, default port : 6379>
    },
    "Host":"<backend host>",
    "FEHost":"<frontend host>"
}
```

The following is sample of config file :
```json
{
    "DB":{
        "Host":"localhost",
        "Port":3306,
        "Username":"root",
        "Password":"root",
        "DbName":"personal"
    },
    "SMTP":{
        "From":"sample@mail.com",
        "Host":"localhost",
        "Port":1025,
        "Username":"username",
        "Password":"password"
    },
    "EmailVerification":{
        "Subject":"Email Verification"
    },
    "ResetPassword":{
        "Subject":"Reset Password"
    },
    "Redis":{
        "Host":"localhost",
        "Password":"",
        "Port":6379
    },
    "Host":"http://personal.localdev.info",
    "FEHost":"http://personal.localdev.info"
}
```

