# Go Interval Notation

This library is used to parse interval notations (e.g. `[v1.2,]` or `[v1,v3)`) for versions in Go. The following format of versions are supported:  
- 1
- 1.2
- 1.2.3
- SemanticVersioning
- Addtionally every version prefixed with a 'v'.

For actually parsing the version and comparing them, it's using [SemVer from @Masterminds](https://github.com/Masterminds/semver)

## Usage

First use `go get` to import the package:
```bash
go get github.com/d0x7/go-interval-notation
```

Afterwards you can use like this:
```go
package main

import "github.com/d0x7/go-interval-notation"

func main() {
  IntervalNotation.InRange("[v1,v2]", "v1.5.3") // == true
}
```

## License

[MIT](LICENSE) © 2023 Dorian Heinrichs
