# Moco Redis Client

Support 3 types of redis server:

- Standalone Redis (Single Redis)

- Sentinel Redis

- Cluster Redis

## Config for Standalone Redis

NOTE: you can let `redis.clientType` is empty string for standalone redis

For json type

```json
{
    "redis": {
        "clientType": "single",
        "address": "localhost:6379",
        "poolSize": 100,
        "password": "",
        "db": 1
    }
}
```

For toml type

```toml
[redis]
clientType = "single"
address = "localhost:6379"
poolSize = 100
password = ""
db = 1
```

## Config for Sentinel Redis

For json type

```json
{
    "redis": {
        "clientType": "sentinel",
        "poolSize": 100,
        "password": "",
        "db": 1,
        "sentinel": {
            "master": "sentinel7000",
            "addresses": [
                "127.0.0.1:7000",
                "127.0.0.1:7001",
                "127.0.0.1:7002"
            ]
        }
    }
}
```

For yaml type

```yaml
redis: 
  clientType: sentinel
  poolSize: 100
  password:
  db: 1
  sentinel: 
    master: sentinel7000
    addresses:
      - 127.0.0.1:7000
```

For toml type

```toml
[redis]
clientType = "sentinel"
poolSize = 10
password = ""
db = 1
[redis.sentinel]
master = "127.0.0.1:30001"
addresses = [ "127.0.0.1:30001", "127.0.0.1:30002", "127.0.0.1:30003"]
```

## Config for Cluster Redis

For json type

```json
{
    "redis": {
        "clientType": "cluster",
        "poolSize": 100,
        "cluster": {
            "addresses": [
                "127.0.0.1:7000",
                "127.0.0.1:7001",
                "127.0.0.1:7002"
            ]
        }
    }
}
```

For yaml type

```yaml
redis: 
  clientType: cluster
  poolSize: 100
  cluster:
    addresses:
      - 127.0.0.1:7000
      - 127.0.0.1:7001
      - 127.0.0.1:7002
```

For toml type

```toml
[redis]
clientType = "cluster"
poolSize = 100
[redis.cluster]
addresses = [ "127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"]
```

