# ShardPilot

## Overview

**ShardPilot** is a query engine designed to route database queries to individual shards based on the approach detailed in [Figma's Database Scaling Blog](https://www.figma.com/blog/how-figmas-databases-team-lived-to-tell-the-scale/). This project provides a proof of concept for handling `SELECT` statements with a Logical Shard ID and other query parameters—a very naive implementation using the OG Postgres query parser's bindings in Golang.

## Features

- **Query Routing**: Routes queries to the correct shard based on Shard ID.
- **Query Transformation**: Converts SQL queries into an Abstract Syntax Tree (AST) and rewrites them to exclude Shard IDs.
- **Database Interaction**: Executes transformed queries on the appropriate shard and returns results.


## Current Functionality (WIP)
<p>
  <img src="/docs/query_engine_poc.png" alt="Query Engine POC" />
  <br>
  <small>My Naive Implementation of this idea from Figma's Blog</small>
</p>

## Setup

Follow these steps to get ShardPilot up and running:

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/harish876/ShardPilot
    ```

2. **Install Dependencies**:

    Navigate to the project root and run:

    ```bash
    go mod tidy
    ```

3. **Setup the Database**:

    - Navigate to the `db/seed` folder:

      ```bash
      cd db/seed
      ```

    - Build and start the Docker containers, then run the seeder script:

      ```bash
      docker-compose up --build
      go run seeder.go
      ```

4. **Start the Server**:

    From the root directory, run:

    ```bash
    go run cmd/main.go
    ```

5. **Query the Endpoint**:

    Send a request to the endpoint with a query containing Shard ID and other columns. For example:

    **Request:**

    ```json
    {
      "query": "SELECT * FROM users WHERE shardId=3 AND name = 'Hannah Harris'"
    }
    ```

    **Response:**

    ```json
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

## TODOs

1. **Topology Management**: Implement and distinguish between logical and physical shards.
2. **Figma Implementation**: Deepen understanding and integration of Figma's approach.
3. **gRPC Conversion**: Convert REST endpoints to gRPC for enhanced performance and scalability.
4. **Physical Planner**: Develop a physical planner with configuration management from a static file.
