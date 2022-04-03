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

## ER Diagram
 Diagram

## Operations
- Insert or update links.
- Iterate over all links.
- Lookup a link by ID.
- Insert of update edges.
- Iterate over all edges. 
- Delete stale edges