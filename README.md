Simple goroutine-safe rate-limiter

Documentation: [https://pkg.go.dev/github.com/Zamony/golimit](https://pkg.go.dev/github.com/Zamony/golimit)

### Example

```go
package main

import (
  "github.com/Zamony/golimit"
)

func main() {
    // Create a new rate-limiter, allowing up-to 10 calls per second
    lim := golimit.New(10, time.Second)
    
    for i := 0; i < 20; i++ {
        if lim.Limit(1) {
            fmt.Println("Over limit!")
        } else {
            fmt.Println("OK")
        }
    }
}
```