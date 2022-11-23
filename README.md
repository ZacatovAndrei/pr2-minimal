# PR2- Minimal

This repository is home to the Lab work #2 from the PR course. It is implemented under the minimal acceptance criteria specifications


## Specifications

`Multiple threads -> 5+`

1. 3 server communicating via HTTP
2. First server is a producer that generates data on multiple threads and sends it to the aggregator
3. Second server is an aggregator. It receives and aggregates data from the producer and sends it to the consumer. The data is extracted from multiple threads
4. Third server is a consumer. It receives data from the aggregator into a shared resource. Extracts it with multiple threads, processes and sends to the aggregator
5. Aggregator aggregates and resends the data to the producer, multple threads, all of that nonsense 