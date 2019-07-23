# Tasks for gyroid

## build

> Builds gyroid

```sh
go build
```

## run

> Runs gyroid

```sh
./gyroid
```

### run go

> Compile and run gyroid main.go file (doesn't build exe)

```sh
go run main.go
```

## install

> Installs gyroid to GOPATH

```sh
go install gyroid
```

## test

> Runs tests

```sh
go test ./...
```

### test junit

> Runs tests and saves them to junit XML format

```sh
mkdir junit && go test -v ./... 2>&1 | go-junit-report > junit/report.xml
```

## clean

> Cleans gyroid

```sh
rm -rf ./gyroid
```
