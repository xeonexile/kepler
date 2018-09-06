# kepler
Go way implementation of CSP paradigm 

You can combine computational pipelines using predefined subset of Springs and Sinks

Spring -> Pipe -> Sink

Each link could be conditional

Each Spring could be connected to multiple Sinks. If there multiple links share the same name

s.linkTo("odd", odd1, OddPredicate)

s.linkTo("odd", odd2, OddPredicate)

then messages will be passed in RoundRobin scenario. 

## Pipe

Pipe is a Spring + Sink composition

There are three types of pipes

1. Transform pipe
2. Broadcaster pipe
3. SelectMany (flatten) roll

go run _examples/text/main.go
