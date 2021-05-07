### Summary

CX executes worker executes the cx program and send back the output. 

The base port number is 9090. 
So for example, when deploying three workers, these are their port numbers:
first worker - 9090
second worker - 9091
third worker - 9092

### Deployment of Worker
For more information, run
```
go run cmd/server.go help 
```

```
go run cmd/server.go -workers=[number of workers to deploy]
```

### Example

For Maze
```
go run cmd/server.go -workers=3
```

### Notes
1. If no arguments are specified, the program will deploy 1 worker at port 9090.
