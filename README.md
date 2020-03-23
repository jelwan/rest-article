# Rest Article

A simple rest api for creating and retrieving articles with tags, the server runs in a docker container and all data is stored in a docker contained mysql database.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for 
development and testing purposes.

### Prerequisites

You will need the following to get started

```
Golang 1.14
Docker 19.03.7
```

### Installing

The following steps will setup the docker database and get the service running

Use the makefile to first create the mysql docker container and wait for the container to start an initalise

```
make docker.database
```

Use the makefile to install the migraton requrements and create the database schema and tabels

```
make db.setup
```

Use the makefile to build the service docker image

```
make docker.buld
```

Use the makefile to run the docker image

```
make docker.run
```

The container should now be running in a docker container with port 8080 exposed on localhost
through the use of a software such as postman or curl you can access the api endpoints on the service

# REST API

The REST API to the rest article is described below.

## Get Article by ID

### Request

`GET /articles/{id}`

    curl -i -H 'Accept: application/json' http://localhost:8080/articles/1

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Mon, 23 Mar 2020 10:36:56 GMT
    Content-Length: 79

    {"id":"1","title":"Get an Article","date":"04-20-2020","body":"Article Body","tags":["tags", "tags"]}

## Create a new Article

### Request

`POST /articles`

    curl -H "Content-Type: application/json"\
    --request POST \
    --data '{\
      "id": "1", \
      "title": "Post an Article", \
      "date": "2020-04-20", \
      "body": "This is how you post an artcile", \
      "tags": ["tags", "tags", "tags"]}' \
      http://localhost:8080/articles


### Response

    HTTP/1.1 201 Created
    Content-Type: application/json
    Date: Mon, 23 Mar 2020 10:39:56 GMT
    Content-Length: 25

    {"success":true,"id":10}

## Get a summary of data about that tag for that day

### Request

`GET /tags/{tagName}/{date}`

    curl -i -H 'Accept: application/json' http://localhost:8080/science/20160922

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Mon, 23 Mar 2020 10:45:17 GMT
    Content-Length: 87

    {"tag":"science","count":1,"articles":["5"],"related_tag":["health","fitness","tech"]}

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

## Authors

* **Andrew Jelwan** - *Initial work* -

