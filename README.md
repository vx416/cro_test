# CRO Test

Crypto dot com backend test. I spent about 6 hrs finishing the project and reviewing the code. I provide a docker-compose to build a server easily and swagger document for testing.

## Contents
- [Getting started](#getting-started)
    - [Run Test](#run-test)
    - [Run Server](#run-server)
- [System Design](#system-design)
    - [Project Layout](#project-layout)
    - [Class Diagram](#class-diagram)
    - [DB Transaction Design](#db-transaction-design)
    - [DB Index Design](#db-index-design)
    - [API Usage](#api-usage)
- [ToDo And Dependencies](#todo-and-dependencies)

### Getting started
```
go mod vendor
```

#### Run Test
**step 1** : build test mysql instance 
```
make run.test.mysql
```

**step 2** : run test
```
go test -p=1 ./internal/test/...
```
note: if test occur `driver: bad connection` error, please run again.

#### Run Server
**run with docker-compose**
[docker-compose install](https://docs.docker.com/compose/install/)
```
make compose.up
```

**run local**
```
go mod vendor
make run.local.mysql
make migrate.local.up
make run.server env=local
```

note: After running server with docker-compose, you can use swagger to test api.

### System Design

#### Project Layout
This project layout is adopted a idea from [golang standard project layout](https://github.com/golang-standards/project-layout) and clean architecture
- `internal/domain` : Domain layer define domain-logic interfaces. It can be used to comply Inversion dependency principle.
- `internal/app` : This package is used to manage dependency and initialize application.
- `internal/model` : Define data structure and its behavior.
- `internal/service`: Servicer layer implement core business logic.
- `internal/delivery` : Handler layer which act as the presenter. It could be RESTful, graphQL or gRPC. Delivery layer will be RESTful API in this project.
- `internal/repository`: Repository provide db operation and encapsulate database-specific 
logic (e.g SQL).
- `internal/test` : Test code and factory object. The test include concurrency testing.
- `pkg/config` : This package used to initialize configuration.
- `build/migrations` : db schema

#### Class Diagram

[![](https://mermaid.ink/img/eyJjb2RlIjoiY2xhc3NEaWFncmFtXG4gICAgY2xhc3MgVXNlcntcbiAgICAgICtTdHJpbmcgRW1haWxcbiAgICAgICtpbnQgSURcbiAgICAgICtpbnQgQ3JlYXRlZEF0aVxuICAgICAgK2ludCBEZWxldGVBdGlcbiAgICB9XG4gICAgY2xhc3MgV2FsbGV0e1xuICAgICAgLWludCBJRFxuICAgICAgLXN0cmluZyBTZXJpYWxOdW1iZXJcbiAgICAgIC1pbnQgVXNlcklEXG4gICAgICAtZGVpbWNsIEFtb3VudFxuICAgICAgK2ludCBDcmVhdGVkQXRpXG4gICAgICAraW50IFVwZGF0ZWRBdGlcbiAgICAgICtpbnQgRGVsZXRlQXRpXG4gICAgICArQ2FuVXNlKHVzZXJJRClcbiAgICAgICtUcmFuc2ZlclRvKHRvV2FsbGV0LCBhbW91bnQpXG4gICAgICArRGVwb3NpdE9yV2l0aGRyYXcodHhLaW5kLCBhbW91bnQpXG4gICAgfVxuICAgIGNsYXNzIFRyYW5zYWN0aW9ue1xuICAgICAgK2ludCBJRFxuICAgICAgK2ludCBLaW5kXG4gICAgICAraW50IEZyb21XYWxsZXRJRFxuICAgICAgK2ludCBUb1dhbGxldElEXG4gICAgICArZGVjaW1hbCBGcm9tV2FsbGV0QmFsYW5jZVxuICAgICAgK2RlY2ltYWwgVG9XYWxsZXRCYWxhbmNlXG4gICAgICArZGVjaW1hbCBUeEFtb3VudFxuICAgICAgK2ludCBDcmVhdGVkQXRpXG4gICAgICArU2V0QmFsYW5jZShmcm9tLCB0bylcbiAgICB9XG4gICAgICAgICAgICAiLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9LCJ1cGRhdGVFZGl0b3IiOmZhbHNlLCJhdXRvU3luYyI6dHJ1ZSwidXBkYXRlRGlhZ3JhbSI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiY2xhc3NEaWFncmFtXG4gICAgY2xhc3MgVXNlcntcbiAgICAgICtTdHJpbmcgRW1haWxcbiAgICAgICtpbnQgSURcbiAgICAgICtpbnQgQ3JlYXRlZEF0aVxuICAgICAgK2ludCBEZWxldGVBdGlcbiAgICB9XG4gICAgY2xhc3MgV2FsbGV0e1xuICAgICAgLWludCBJRFxuICAgICAgLXN0cmluZyBTZXJpYWxOdW1iZXJcbiAgICAgIC1pbnQgVXNlcklEXG4gICAgICAtZGVpbWNsIEFtb3VudFxuICAgICAgK2ludCBDcmVhdGVkQXRpXG4gICAgICAraW50IFVwZGF0ZWRBdGlcbiAgICAgICtpbnQgRGVsZXRlQXRpXG4gICAgICArQ2FuVXNlKHVzZXJJRClcbiAgICAgICtUcmFuc2ZlclRvKHRvV2FsbGV0LCBhbW91bnQpXG4gICAgICArRGVwb3NpdE9yV2l0aGRyYXcodHhLaW5kLCBhbW91bnQpXG4gICAgfVxuICAgIGNsYXNzIFRyYW5zYWN0aW9ue1xuICAgICAgK2ludCBJRFxuICAgICAgK2ludCBLaW5kXG4gICAgICAraW50IEZyb21XYWxsZXRJRFxuICAgICAgK2ludCBUb1dhbGxldElEXG4gICAgICArZGVjaW1hbCBGcm9tV2FsbGV0QmFsYW5jZVxuICAgICAgK2RlY2ltYWwgVG9XYWxsZXRCYWxhbmNlXG4gICAgICArZGVjaW1hbCBUeEFtb3VudFxuICAgICAgK2ludCBDcmVhdGVkQXRpXG4gICAgICArU2V0QmFsYW5jZShmcm9tLCB0bylcbiAgICB9XG4gICAgICAgICAgICAiLCJtZXJtYWlkIjoie1xuICBcInRoZW1lXCI6IFwiZGVmYXVsdFwiXG59IiwidXBkYXRlRWRpdG9yIjpmYWxzZSwiYXV0b1N5bmMiOnRydWUsInVwZGF0ZURpYWdyYW0iOmZhbHNlfQ)

```
|users| -has-many-> |wallets| -has-many-> |transactions|
```

#### DB Transaction Design
Transferring money between wallets should consider race condition when multiple requests are processed simultaneously.

The most safe way is used to `FOR UPDATE` SQL to put a write lock on record. However, we can use optimistic lock to have better performance.

Therefore, I use following SQL to update wallet balance.
```sql
UPDATE wallets SET amount=amount-? WHERE amount >= ? AND serial_number = ?;
```
If this update statement show no rows affected, you can determine amount of this wallet is insufficient.


#### DB Index Design
Considering that users might need to review their transaction, we need to add proper index on `transactions` table.

The most common search scenario will be search transactions during a period, so I create two composite indexes `idx_to_wallet_id_created_ati` and `idx_from_wallet_id_created_ati` for this search scenario.

These indexes does not include `kind` column because kind column can be multiple values in search condition.


#### API Usage
**swagger URI** : `http://localhost:8080/swagger/index.html`

**Authentication** : 
After creating an account, you can use `/api/v1/auth` api to get your token, then put your token in `Authentication` header with `Bearer` schema.
`POST /api/v1/signup` : create an account and wallet
`POST /api/v1/auth` : get token


**Wallet** : 
`POST /api/v1/wallets` : create wallet
`GET /api/v1/wallets/:serial` : get wallet and its transactions in 7dys.
`GET /api/v1/wallets` : get all wallets.

**Transaction** : 
`POST /api/v1/transfer` : transfer money
`POST /api/v1/deposit` : deposit money
`POST /api/v1/withdraw` : withdraw money 

### ToDo And Dependencies

#### ToDo
- List wallets and transactions meanwhile solve n+1 query problem.
- Allow users to search their transactions.
- Provide delete api by implementing soft delete operation.
- Add cicd flow (e.g travis-ci, ansible).
- Add kubernetes deployment and Helm chart.
- Add redis rate limiter (e.g leaky bucket).
- Improve security with CORS, HTTPs and CSRF secret.
- Add performance monitoring mechanism (e.g Elastic APM and OpenTelemetry).
- Provide gRPC endpoints.

#### Dependencies
- [echo](https://github.com/labstack/echo) for http server.
- [fx](https://github.com/uber-go/fx) for dependency injection.
- [testify](https://github.com/stretchr/testify) for building test suite.
- [gogo-factory](https://github.com/vx416/gogo-factory) for building factory object.
- [sqlxx](https://github.com/vx416/sqlxx) for extension of [sqlx](https://github.com/jmoiron/sqlx).
- [viper](https://github.com/spf13/viper) for configuration management.
- [zap](https://github.com/uber-go/zap) for logger.
- [jwt-go](https://github.com/dgrijalva/jwt-go) for JWT token.
- [goose](https://github.com/pressly/goose) for db migrations.
