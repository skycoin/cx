# `cxgo` Endpoints

`cx` serves an HTTP interface when it is run with the `--web` flag specified and supports the following endpoints:

```bash
# Query program meta data.
$ curl http://localhost:5336/program/meta
{
  "used_heap_memory": 0,
  "free_heap_memory": 0,
  "stack_size": 0,
  "call_stack_size": 0
}

# Query program packages.
$ curl http://localhost:5336/program/packages
[
  "http",
  "cipher",
  "regexp"
]

# Query a specified package.
$ curl http://localhost:5336/program/packages/http
{
  "functions": [],
  "structs": [
    {
      "name": "URL",
      "signature": "URL struct { Scheme str; Opaque str; Host str; Path str; RawPath str; ForceQuery bool; RawQuery str; Fragment str; Close bool; Uncompressed bool; }",
      "type": 15,
      "type_name": "struct"
    },
    {
      "name": "Request",
      "signature": "Request struct { Method str; URL *URL; Header [][]str; Body str; }",
      "type": 15,
      "type_name": "struct"
    },
    {
      "name": "Response",
      "signature": "Response struct { Status str; StatusCode i32; Proto str; ProtoMajor i32; ProtoMinor i32; Body str; ContentLength i64; TransferEncoding []str; }",
      "type": 15,
      "type_name": "struct"
    }
  ],
  "globals": []
}```