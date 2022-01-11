# Golang authentication service

## Local development
1. Configure `env` file
2. Run postgres and redis in docker container
```shell
$ make compose-up
```
3. Build and run the app
```shell
$ make run
```

Make sure Makefile includes correct `env` file.

## Links
- [Evrone Go clean template](https://github.com/evrone/go-clean-template)
