# dating-api

## Running the service
The service can be run using the provided docker-compose file: 
```
docker-compose up --build -d
```

## Documentation
The documentation for the api can be viewed at:
```
http://localhost:8080/dating-api/v1/swagger/index.html#/
```

## Authentication
Once the user has logged in, they must use the signed JWT for all other requests. 
```
Authorization: Bearer <YOUR_JWT>
```
By default the JWT will expire after 5 minutes after which you must request a new one. The authentication for each request is handled through custom middleware defined in `router.go`. This validates the JWT, and sets the requesting userID in the context to allow the usecases to access it.

## Assumptions
- Gender filtering: that being able specify multiple genders is an appropriate option for the target customers.