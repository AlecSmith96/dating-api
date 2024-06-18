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