
# How To Run


## Prerequisites

- Install [Go](https://golang.org/doc/install)


- Install [gomock](https://github.com/golang/mock)


- Install [go-migrate](https://github.com/DavidHuie/gomigrate/blob/master/README.md)

- PostgreSQL

    Follow [PostgreSQL download](https://www.postgresql.org/download/).
    
## Run Instruction    

### Using Docker
 ```
    make docker-compose
```

### Manual
- copy env config from env.example
 ```
    cp env.example .env
```

- set env variables with prefix POSTGRES_ according to your Postgres settings

- download depedencies
```
    make tidy
```

- run 
```
    make run
```

## How To Test
run this command for running test with coverage
```
    make test-coverage
```
