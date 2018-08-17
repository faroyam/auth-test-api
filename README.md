# auth-test-api

Simple REST API with JWT authentication written in GO.

### Installing

```
go get github.com/faroyam/auth-test-api
```

## Getting Started

1. Generate SSL keys for HTTPS and JWT 
```
cd keys && bash keys.sh
```
2. Run app
```
cd .. && go build app.go && ./app
```