# errors package

- every error should be wrapped with message by passing format like so:
```go
if err := someService(); err != nil {
    log.Fatal(errors.Newf("someService: %w", err) // will create NoType error)
    log.Fatal(errors.BadRequest.Newf("someService: %w", err) // will create error with BadRequest type
}
```

