## Shard Pilot

1. Implemented a simple query engine to route database queries to individual shards based on the article (https://www.figma.com/blog/how-figmas-databases-team-lived-to-tell-the-scale/)
2. A simple proof of concept solution was done for select statements of Shard ID and other columns to be queried.   


# Setup
1. Clone the repository ```git clone https://github.com/harish876/ShardPilot```
2. Install Deps ```go mod tidy```
3. Navigate to the db/seed folder and run the docker-compose file to set up the database and run the seeder scripts like so ``` cd db/seed && docker-compose up --build && go run seeder. go```
4. Run ``` go run cmd/main.go``` from the root directory to start the server.
5. Hit the endpoint with a query that contains Shard ID and other columns queried
   Request:
    ```
       {
        "query":"select * from users where shardId=3 and name = 'Hannah Harris'"
       }
     ```
   Response: The query engine converts the query to an AST, gets the Shard ID, gets the correct connection string for that shard rewrites the query without the Shard ID, and gets the rows from the table.
     ```
        {
       "message": "Query on ShardID 3. Query Executed on the DB: SELECT * FROM users WHERE name = 'Hannah Harris'",
       "data": [
        {
            "UserID": 34,
            "Name": "Hannah Harris",
            "PhoneNumber": "+71-5781028656",
            "Email": "hannah_harris@example.com"
        }
       ],
       "error": ""
        }
     ```

**TODOS**
1. Topology management and the distinction between logical shards and physical shards.
2. Improving / Learning more about the Figma implementation in detail.
3. Converting REST Endpoints to GRPC because why not?
4. Physical Planner with proper configuration management from a static file.
