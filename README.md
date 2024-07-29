## Shard Pilot

1. Implementation of a query engine to route database queries to individual shards. Based on the ideas proposed by this article (Figma Postgres Sharding)[https://www.figma.com/blog/how-figmas-databases-team-lived-to-tell-the-scale/]


2. Includes Lightweight logical query planner and physical query planner.
    - Logical query planner gets the query type (SELECT,INSERT,UPDATE,DELETE) and the shardId
    - Physical query planner maps the logical shardIds to physical shards (**TODO**)

3. Current HTTP layer implementation uses REST. Plan to convert to GRPC/Cap n Proto

4. Topology management for Physical query planner, Docker Compose to set up database shard (**TODOS**) 
