# 🔬 chi-ratelimit
> *Simple production-ready ratelimiter for Chi applications*

## What is this?
**chi-ratelimit** is middleware to implement production-safe ratelimiting into your [go-chi](https://github.com/go-chi/chi) applications.

## Example
```go
package main

import (
  "github.com/go-chi/chi/v5"
  "github.com/Noelware/chi-ratelimit"
  "github.com/Noelware/chi-ratelimit/redis"
  "fmt"
)

func main() {
  ratelimiter := ratelimit.NewRatelimiter(
    ratelimit.WithStorage(redis.New().WithClient(<redis client here>)),
    ratelimit.WithDefaultTime(func() time.Time { time.Now().Add(1 * time.Hour) }),
  )
  
  router := chi.NewRouter()
  router.Use(ratelimiter)
  
  http.ListenAndServe(":3000", router)
}
```

## License
**chi-ratelimit** is released under the **MIT License** by Noelware.
