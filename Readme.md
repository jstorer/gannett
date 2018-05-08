
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

## Install
The latest version of the API is pushed to dockerhub as well as hosted in a kubernetes cluster.
To use the docker version run the commands

first get the latest version of the docker image

```docker pull jstorer/gannett```

next run the image

```docker run -p 8080:8080 jstorer/gannett```

you will see "*...Supermarket Server Starting...*" if everything worked as intended.

## End Points








