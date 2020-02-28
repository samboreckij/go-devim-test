# ENV

ENV is utility functions for environment-based configs.

```go

env.Get("POSTGRES_PASSWORD").String("passw0rd")
env.Get("POSTGRES_MAX_CONNECTIONS").Int(10)
env.Get("POSTGRES_SSL_ENABLED").Bool(true)
env.Get("POSTGRES_PING_TIMEOUT").Duration(30 * time.Second)

```
