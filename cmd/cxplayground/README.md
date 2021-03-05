# Instructions

- Run `go run main.go` inside this directory. You'll see the following output:

```
Starting web service for CX playground on http://127.0.0.1:5336/
```

- In another terminal, run this as a hello world example:

```
curl -H "Content-Type: application/json" -X POST -d '{"code": "package main; func main() {str.print(\"Hello, world!\");}"}' http://localhost:5336/eval
```

- Get playgournd page
```
curl -H "Content-Type: application/json" -X GET http://localhost:5336/playground
```

- Get example file list
```
curl -H "Content-Type: application/json" -X GET http://localhost:5336/playground/examples
```

- Get example file content
```
curl -H "Content-Type: application/json" -X POST  -d '{"examplename": "example1.cx"}' http://localhost:5336/playground/examples/code
```