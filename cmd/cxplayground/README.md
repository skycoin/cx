# Instructions

- Run `go run main.go` inside this directory. You'll see the following output:

```
Starting web service for CX playground on http://127.0.0.1:5336/
```

- In another terminal, run this as a hello world example:

```
curl -H "Content-Type: application/json" -X POST -d '{"code": "package main; func main() {str.print(\"Hello, world!\");}"}' http://localhost:5336/eval
```
