# 12 factor app demo

Ce repository est utilisé dans le cadre de la démo des 12 factors app

## Lien

https://12factor.net

## List of variables

TWELVE_MAIL_PASSWORD       
TWELVE_MAIL_FROM    
TWELVE_SERVER_ADDRESS
TWELVE_DB_MONGO_ADDRESS
TWELVE_AUTH_SECRET
TWELVE_AUTH_PASSWORD

## Quick start with docker

### Run DB in docker


```bash
cd db
docker compose up
```

### Run main app
in a second bash

```bash
docker build -t mail:latest .
docker run --rm --name mail -e "TWELVE_MAIL_PASSWORD=<insert your mail password>" -e "TWELVE_MAIL_FROM=<insert your mail address>" -e "TWELVE_SERVER_ADDRESS=<insert your smtp address>" -e "TWELVE_DB_MONGO_ADDRESS=mongodb://localhost:27017" -e "TWELVE_AUTH_SECRET=thisisnotastrongsecret" -e "TWELVE_AUTH_PASSWORD=banana" -p 8080:8080 mail:latest
```