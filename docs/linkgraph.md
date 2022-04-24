# Link Graph
## Link
- Web pages which are either processed or just discovered.
```
ID 
URL
RetrievedAt
```
## Edge
- Represents a uni-directional connection between two links in the graph.
- This will give information about origination of a link.
```
ID
SourceLink
DestinationLink
UpdatedAt
```
## Operations
- Insert or update links.
- Iterate over all links.
- Lookup a link by ID.
- Insert of update edges.
- Iterate over all edges.
- Delete stale edges

## Database
- CockroachDB
    - CockroachDB can easily scale horizontally, by just increasing number of nodes.
    - ACID compliant
    - Support distributed SQL transactions
    - SQL is compatible with PostgreSQL syntax and hence we can use postgres drivers.

## Migrations
- `golang-migrate/migrate` will be used to manage schema migrations
- it ensures that migrations are run only once by maintaining state in two additional tables.