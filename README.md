# Nolang

A basic template for golang web apps. It's a work in progress

## App Dependencies

- [chi](https://github.com/go-chi/chi) as the router
- [zerolog](https://github.com/rs/zerolog) for logs
- [envconfig](https://github.com/kelseyhightower/envconfig) and [godotenv](https://github.com/joho/godotenv) for configuration
- [context](https://golang.org/pkg/context/) for managing request contexts

## Setup

Make sure you're using `golang >= 1.13`(due to dependency on `go mod`)

Install [air](https://github.com/cosmtrek/air)

```sh
# Do this outside the project so you can use it with other project
go get -u github.com/cosmtrek/air
```

Run air

```sh
air
```

## Structure

- Everything related to configuration will go to the `config` dir
- Your handlers should go to the `controller` dir.

**Note**: This is a `Work-in-Progress`, expect constant changes
