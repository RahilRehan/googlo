package graphtest

/*
This package contains following group of tests:
- Link/Edge upsert tests
- Concurrent Link/Edge iterator
	to verify no data races occur when multiple iterator instances are present.
- Partitioned iterator tests
	to verify that if we split our graph into N partitions and each partition has an iterator
		- verify no link/edge is in more than one partition
		- all iterators together will process all links/edges
		- edge iterators ensure that edge is present in same partition as source link
- Link lookup tests
- Stale edges removal tests
*/

