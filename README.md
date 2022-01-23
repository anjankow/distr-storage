# cloud-storage-system

Cloud storage system for key-value pairs in hash-based structure. The storage is distributed, yet seems to be one DB from the user perspective.

## Architecture

The storage system consists of a data-access-api application and arbitrary number of nodes. The API application is independent from nodes and can't access their DB directly.

### data-access-api
Is an entry point to the system. It serves as a directory service for the nodes. Using a hash function and modulo operation it assignes each key to a bucket. On a request reception, it forwards it to the corresponding bucket with an HTTP request.

#### API
 - POST `/config`
 
 Request body:
 ```json
 {
	"collection": "myCollection",
	"nodes": [
		"node0:8080",
		"node1:8080",
		"node2:8080"
	]
}
```

Response: status code

The config request is used to configure the API application in runtime. `nodes` param informs about the addresses of all the nodes in the network. The application sends a healtcheck to each of the nodes, waiting until they are responsive.

- PUT `/doc`

Request body:
```json
{
	"key": "favouriteMuffin",
	"value": {
		"name": "brownie",
		"topping": "chockolate"
    }
}
```

- GET `/doc`

    filters:

    - `key=favouriteMuffin`
    gets a single object from with the key "favouriteMuffin"

        Response body:
    ```json
    {
      "node": "node0:8080",
        "value": {
            "key": "favouriteMuffin",
            "value": {
                "name": "brownie",
                "topping": "chockolate"
            }
        }
    }
    ```
    The reponse appends information about the node holding the given key.


    - `all=true`
    gets all the objects belonging to the preconfigured collection from all the nodes in the network

        Response body:
    ```json
    {
        "node_data": [
            {
                "node": "node0:8080",
                "data": [
                    {
                        "id": "favouriteMuffin",
                        "timestamp": "2022-01-22T18:12:32.532Z",
                        "content": {
                            "name": "brownie",
                            "topping": "chockolate"
                        }
                    }
                ]
            },
            {
                "node": "node1:8080",
                "data": [
                    {
                        "id": "whiteMuffin",
                        "timestamp": "2022-01-22T18:12:32.532Z",
                        "content": {
                            "name": "brownie",
                            "topping": "frosting"
                        }
                    },
                    {
                        "id": "redMuffin",
                        "timestamp": "2022-01-22T18:12:32.532Z",
                        "content": {
                            "name": "brownie",
                            "topping": "red glaze"
                        }
                    }
                ]
            }
        ]
    }
    ```

 - DELETE `/doc?key=redMuffin`
    
    Response: status code

### nodes
Each node application is managing one bucket. For buckets mongo DB is used, as an efficient key-value database.

A node is configured via the environmental variables: URI of the DB and DB name.

## deployment

The service is started throgh executing:
```
./start_db.sh <number of buckets>
```

The script creates a docker-compose file for the desired number of buckets and starts the system.

There are 3 network defined in the docker-compose file:
 - external - which connects the data-access-api with the external world (ideally)
 - api-network - between nodes and the api application
 - nodeX-network - between a node and its mongo DB instance

This layout allows for the following separation:
 - only node applications can access its mongo DB, they can't access the DBs belonging to another nodes
 - a client can't access the node applications directly, only through the api app

## client

The client package holds the implementation of a simple client and a tester class. 
The tester is using the client to execute a test sequence and generate the report.

To start the test:
```sh
cd client/
task run
```
