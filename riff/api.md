# API
## nodes
> post /api
```graphql
{
    nodes {
        name
        ip
        port
        dataCenter
        state
        isSelf
        version
    }
}
```

```json
{
    "data": {
        "nodes": [
            {
                "dataCenter": "dc1",
                "ip": "192.168.1.220",
                "isSelf": true,
                "name": "node1",
                "port": 8630,
                "state": "Alive",
                "version": 0
            }
        ]
    }
}
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
