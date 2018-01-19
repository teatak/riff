# API
## nodes
> get /api/nodes

```json
[
    {
        "name": "node1",
        "dataCenter": "dc1",
        "ip": "192.168.3.2",
        "port": 8630,
        "state": 1,
        "snapShot": "c18753f2a0be4531bf2bfc68327a98d6180adfa6"
    },
    {
        "name": "node2",
        "dataCenter": "dc1",
        "ip": "192.168.3.2",
        "port": 8631,
        "state": 1,
        "snapShot": "59f255e3ca0af4d45609804cfd52b11a32c6bf34"
    }
]
```

## node
> get /api/node/:name

```json
{
    "name": "node1",
    "dataCenter": "dc1",
    "ip": "192.168.3.2",
    "port": 8630,
    "state": 1,
    "snapShot": "c18753f2a0be4531bf2bfc68327a98d6180adfa6",
    "services": [
        {
            "name": "mongod",
            "ip": "192.168.3.2",
            "port": 27017,
            "state": 3
        },
        {
            "name": "ping",
            "state": 1
        }
    ]
}
```

## services

> get /api/services

```json
[
    {
        "name": "mongod"
    },
    {
        "name": "ping"
    }
]
```

## sevice

> get /api/service/:name -> /api/service/:name/all

> get /api/service/:name/:command

command 
* alive
* all

```json
{
    "name": "mongod",
    "nodes": [
        {
            "name": "node1",
            "dataCenter": "dc1",
            "ip": "192.168.3.2",
            "port": 27017,
            "state": 1
        },
        {
            "name": "node2",
            "dataCenter": "dc1",
            "ip": "192.168.3.3",
            "port": 27017,
            "state": 1
        }
    ]
}
```

> post /api/service/:name/:command

command:
* start
* stop
* restart

```json
{
    "status": 201
}
```

status:
* 200 success
* 201 success
* 400 command missing
* 404 not found
