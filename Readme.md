
# Supermarket REST API

## Overview
This is an API for use at a local
grocery store that can add, delete, update,
and fetch all produce in the system.
It is set up to be continuously
delivered using **Golang, Docker,
Kubernetes, and Travis-CI.**

The produce database is held in a single,
in memory array of data. Each item in the
database includes **name, produce code,
and unit price.** The produce code is alphanumeric
and case insensitive with the format of
`XXXX-XXXX-XXXX-XXXX` where *X* is any number or letter.
The unit price is a number with up to two
decimal places. The name is alphanumeric.

The API was designed with RESTful
principles in mind, containing proper
response codes and HTTP methods. Tests
were written to accompany the API to
help ensure correctness.

## Install & Run
### Install and Run Docker Image
The latest version of the API is pushed to dockerhub as well as hosted in a kubernetes cluster.
To use the docker version run the commands

first get the latest version of the docker image

```docker pull jstorer/gannett```

next run the image

```docker run -p 8080:8080 jstorer/gannett```

you will see "*...Supermarket Server Starting...*"
if everything worked as intended then port 8080 of the container
should be bound to port 8080 of the local machine.
This will result in the following end points using `http://localhost:8080{end point}`

## End Points

### GET Method
#### Get All Items
`/api/produce`

This method will return all produce items in the database in JSON.
#### Get One Item
`/api/produce/{produce_code}`

This method will return the produce item of the given *{produce_code}* in JSON or an error if it does not exist.

### POST Method
#### Create New Item
`/api/produce`

This method will accept a JSON request from the body and add it to the database assuming all neccesary fields are populated correctly.

*Input Field Requirements*:
* Produce Code - required, unique
* Name - required
* Unit Price - required

An error will be returned if requirements are not met. A JSON response of the fields will be returned upon success.
#### Update Existing Item
`/api/produce/{produce_code}`

This method will accept a JSON request from the body and add it to the database assuming all neccesary fields are populated correctly.

*Input Field Requirements*:
* Produce Code - required, unique
* Name - required
* Unit Price - required

An error will be returned if requirements are not met or *{produce_code}* does not exist. A JSON response of the now updated item will be returned upon success.

### DELETE Method
#### Delete Existing Item
`/api/produce/{produce_code}`

This method will delete the item from the database that contains *{produce_code}*. An error will be returned if *{produce_code}* does not exist. Upon success a JSON response containing the deleted item will be returned.

## Application Structure
The application contains several major parts. The API itself was written in Go which can be containerized and pushed to dockerhub at [https://hub.docker.com/r/jstorer/gannett/] where it can be ran as dicussed earlier. It is also hosted on a Kubernetes cluster via an image that is pushed to the Google Container Registry(GCR).
### API (Go)
The API is broken up into several files for the sake of organization and clarity.
#### File Summary
##### main.go
This file functions as a kick off point to start the api package and initialize the database and start the server listening.
##### handlers.go
This is where the routing is set for the different end points that were referenced earlier.
##### api.go
The handler functions from the routing and a few helper functions are contained inside.
##### model.go
The data structures and their methods are contained here along with functions that directly manipulate the database.
##### api_test.go
Tests to ensure API is working correctly are contained inside of here

#### General Structure
The in memory array structure to store data,named `DBObject`, is a struct that holds a mutex, which will be used to help prevent race conditions, and a type `ProduceItem` slice.
```
type ProduceItem struct {
    ProduceCode string `json:"produce_code"`
    Name        string `json:"name"`
    UnitPrice   string `json:"unit_price"`
}

type DBObject struct {
    mu   sync.RWMutex
    Data []ProduceItem
}
```
Upon starting, the application will select either a production or test database via a flag in `api.Initializer(isTesting bool)`, so testing and running can have their own data sources. Then the routes will be set as, seen in *handlers.go*, and the application will begin listening on port 8080. Depending on the request one of the handler functions will fire:

##### Handler Functions
These are the functions set by the router to handle incoming requests.

###### handleGetAllProduce(ResponseWrite, *Request)
This function sends a request to the database to fetch all produce items through a goroutine,`getAllProduceItems(chan []ProduceItem)`,
then returns them on a channel then finally returns them in JSON format with a 200 status code.

###### handleGetProduceItem(ResponseWrite, *Request)
This function first retrieves the produce code from the URL and
determines if it is valid. If it is not valid it triggers a status 400 error.
If it is valid it fires a goroutine,`getProduceItem(chan ProduceItem)`, to fetch that particular item and waits
for a response via a channel. If the database returned an item it is displayed in JSON
 along with a 200 status code. If it is not found a 404 status code is triggered.

###### handleCreateProduceItem(ResponseWrite, *Request)
This function first parses the JSON body request into a `ProduceItem` type then
checks to see that all fields are valid and filled in by calling the `ProduceItem`
method `validateProduceItem()`. If validation fails a status code 400 is triggered
along with a JSON response of the errors. If the `ProduceItem` is valid a goroutine, `createProduceItem(ProduceItem,chan ProduceItem)`,
is triggered to create an item with the data passed back through a channel.
 If the produce code already exists in the data a status code 409 is triggered
 if not a 201 status code is triggered with the JSON of the `ProduceItem` returned.

###### handleUpdateProduceItem(ResponseWrite, *Request)
This function checks if the produce code passed in from the URL is valid,
if it is not a status code 400 is triggered. If it is the JSON from the request
body is placed into a `ProduceItem`. This is then validated the same way as the
create function. Upon validation success a go routine,
`updateProduceItem(produce_code string, ProduceItem, chan ProduceItem)`, is called and
passes the updated item back through a channel. If the produce code was not found a
status code 404 is triggered or if the changed produce code already exists a status
409 is triggered. Otherwise a status 200 is triggered and the updated item contents
are returned as a JSON.

##### Handler Helper Functions
These functions are used to perform some frequent duties inside the handler functions

###### func isValidProduceCode(produceCode string) bool
This function accepts a produce string and validates via the regex expression `^[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}$`
to determine if it is valid or not and returns true if valid or false if not. This expression checks that code is four groups of four alphanumeric characters.

###### func isValidUnitPrice(unitPrice string) bool
This function accepts a unit price string and validates via the regex expression `^\$(([1-9]\d{0,2}(,\d{3})*)|(([1-9]\d*)?\d))(\.\d\d?)?$`
which requires a dollar sign followed by numbers with or without correct comma seperation but not incorrect comma seperation and at most 2 trailing decimals. If valid returns true and if not valid returns false.

###### func isValidName(name string) bool
This function accepts a name string and validates via the regex expression `\w+(?: \w+)*$`
which will allow no whitespace before or trailing whitespace after a single space set of alphanumeric characters.
If valid returns true if not valid returns false.

###### func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{})
This function accepts a ResponseWriter, status code, and payload and then JSON encodes it via json.Marshall()
and then writes the corresponding body and headers in JSON format to be displayed.

###### handleDeleteProduceItem(ResponseWrite, *Request)
This function first checks if the produce code passed in from the URL is valid,
if it is not a status code 400 is triggered. If the produce code is valid
a goroutine,`deleteProduceItem(ProductionItem, chan ProduceItem)`, is triggered and passes
the produce item back through a channel. If the code was not found a status 404
is triggered, if it was found a status 200 is triggered and the deleted produce item
is returned as a JSON.

##### Model Methods (Go Routines)
These are the functions that change values in the database or are methods of created data types.

###### getAllProduceItems(chan)
`RLock()`s the database and returns all produce items on channel then `RUnlock()`s the database.

###### getProduceItem(string, chan ProduceItem)
`RLock()`s the database then searches for produce code. If the code
is found it returns the corresponding item on a channel and if not
found returns an empty item on a channel then `RUnlock()`s the database.

###### createProduceItem(ProduceItem, chan ProduceItem)
`Lock()`s the database and brings the produce code to upper case since
it is case insensitive and will give consistency to how the data is presented.
If the code already exists an empty ProduceItem is returned on the channel. Othewise,
The data is appeneded to the database and the created item is returned on the channel.
The database is then `Unlock()`ed at the end of either case.

###### updateProduceItem(string, ProduceItem, chan ProduceItem)
`Lock()`s the database and brings the produce code to upper case since
it is case insensitive and will give consistency to how the data is presented.
It then checks to see if the produce code to be updated exists. If it does exist it checks
if the new value already exists and returns a ProduceItem with a code of '0' if it does on a channel.
Otherwise it changes the values of the database at the found location with the new information and returns the updated item on the channel.
If the item to be updated is not found an empty ProduceItem is returned on the channel. At the end
of any case the database is `Unlock()`ed.

###### deleteProduceItem(ProduceItem, chan ProduceItem)
`Lock()`s the database and searches for the produce code given. If the code is found
that item is removed from the database and its information retruend on the channel.
If it is not found an empty ProduceItem is returned. At the end of either case
the database is `Unlock()`ed.

##### Testing
Test code is located in api_test.go and done in table format with assistance from the testify package [https://github.com/stretchr/testify]
to faciliate easy to read and write test code.

### Docker - Multi-stage build
The docker build specifics are located in the Dockerfile. It is a multi-stage docker
build to keep the size of the image down. It first uses the golang image to build
the app by copying needed files into the correct places and then getting any needed
dependencies and finally building the app file.

The next stage of the build then takes over and copies the built application file from
the previous stage and executes it.
### Kubernetes
A kubernetes cluster was used as a place to deploy the application to. Kubernetes
uses a docker image that is pushed to the google container registry. This image
is then placed in a container cluster where a Pod(s) contains the application.
It is then exposed to the Internet by creating creating a service resource through
*kubectl* which provides networking and IP support by creating an external IP
and Load Balancer.

It is possible to scale up the application when needed by adding replicas to the deployment
resource using *kubectl scale*.

### Travis-CI
Travis-CI is used for continuous integration via a travis.yml file and github
association. Every time the application is commited to Github Travis-CI:
* builds the golang code
* runs tests
* builds and pushes an image to dockerhub
* builds and pushes an image to gcr which is then used to update the kubernetes cluster by executing
the *deploy-production.sh* script.
* notifies users of success or failure.










