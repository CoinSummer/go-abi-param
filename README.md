## Golang parse abi params On EVM

---

### Features
1. parse evm abi argument type
2. When parsing data types and numbers, this library will convert all parameter values to string type and ignore whitespace. It will ultimately parse and return the corresponding Go variable type.

### Usage
```go
package main

import (
	"fmt"

	ap "github.com/CoinSummer/go-abi-param"
)

func main() {
    param, err := ap.NewAbiParam("uint8[]", "[0,1]")
    if err != nil {
        panic(err)
    }
    res, err := param.Parse()
    if err != nil {
        panic(err)
    }
    fmt.Println(res) // []uint8{0,1}
}
```