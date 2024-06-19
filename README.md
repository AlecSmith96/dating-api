# dating-api

A simple REST server for users to login to the platform, view other users based on criteria, and swipe on them with their preference.

## Service design and architecture
The service is implemented using [Clean Architecture](https://celepbeyza.medium.com/introduction-to-clean-architecture-acf25ffe0310) in order to make it as maintainable and testable as possible. It does 
this by decoupling the application code from the business logic, allowing for changes to the application code to be made
without impacting the business logic (or vice versa). It separates the code into four distinct layers, each responsible 
for different things:
- `drivers`: This layer is responsible for delivering data to and from the usecases. In this service, that is done through a [gin](https://github.com/gin-gonic/gin) router.
- `adapters`: This layer is responsible for the implementation details of any interaction with external resources, such as databases, APIs or other services. They implement interfaces in the `usecase` layer to allow the implementation to be swapped out if necessary.
- `usecases`: This layer is responsible for the business logic of each endpoint. It depends on interfaces to decouple the implementation details from the business logic.
- `entities`: This layer stores any internal service structs, it is used by all other layers to parse data across them.

## Running the service
The service can be run using the provided docker-compose file: 
```
docker-compose up --build -d
```
For convenience, a test user gets created when the service starts up, it can be logged in using the following credentials:
```
{
  "email": "admin",
  "password": "admin"
}
```

## Documentation
The documentation is generated from the code using [swagger](https://github.com/swaggo/gin-swagger), it can be viewed at:
```
http://localhost:8080/dating-api/v1/swagger/index.html#/
```

## Authentication
Once the user has logged in, they must use the signed JWT for all other requests. 
```
Authorization: Bearer <YOUR_JWT>
```
By default the JWT will expire after 5 minutes after which you must request a new one. The authentication for each request is handled through custom middleware defined in `router.go`. This validates the JWT, and sets the requesting userID in the context to allow the usecases to access it.

## Running the tests
Due to time constraints 100% test coverage couldn't be achieved, but tests for each layer were written. You can run the tests using the following command:
```
make test
```
The tests require the `postgres` container to be running due to the migration tests running each migration to test for issues.
There are 3 types of tests in this service:
- Migration tests: Tests each migration in isolation to ensure they work as intended. Each test creates a new database to run the migration on and tests the data is correct when entered.
- Adapter unit tests: Unit tests for the postgres adapter, using mocked SQL calls to test the interaction with postgres and the returning data to the usecase.
- Usecase behavioural tests: Unit tests that specifically test the business logic of each endpoint, using mocked calls to the adapter layer.

In addition to these tests, integration tests would be written to test the interactions between this service and any others, against real resources in a realistic environment.

## Bonus Task
Due to time constraints, I have been unable to implement the bonus task of sorting profiles by attractiveness. 
My plan for implementation of this was to base the attractiveness on the users match success percentage. To get this you 
calculate the percentage of matches that a user gets for their swipes. If a user swipes 'YES' 100 times and gets 10 matches, 
they would be marked as more 'attractive' than a user who swipes 'YES' 50 times but gets 3 matches. Then, I would write 
an SQL query to get this figure for each user and use this to sort the returned slice of users in the `/discover` endpoint. 


## Assumptions
During the implementation of the task, I made the following assumptions: 
- I made the assumption that the target users of the app would want the functionality to select multiple genders when 
searching for matches. Depending the target audience, this might not be appropriate.
- I also made the assumption that encrypting the user passwords using a standard encryption algorithm would not be 
necessary for this task, as it was not explicitly mentioned and wouldn't display any particular skills due to using
standard algorithms.
- Finally, I also made the assumption that pagination was not necessary for the `/discover` endpoint. 
