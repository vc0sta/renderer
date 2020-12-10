# Endpointer

This application was created to test a wide range of endpoints,
from http to databases, it can test a single time or either keep watching
until the target application respond succesfully.

```sh
> go run main.go check http -h
```

```text
Check if a http/https is reachable.
It will not test if the returning code is 2XX.

Usage:
  endpointer check http <url> [flags]

Examples:
endpointer check http https://google.com

Flags:
  -h, --help          help for http
      --timeout int   how many seconds should a watch run (default 3600)
      --watch         keep watching command, retries connection each 2s.
```
# renderer
