# netology-slowly

Example project for Netology Golang expert position

## Requirements

* Go
* make (in order to use Makefile)
* golangci-lint (in order to use linters)

### Usage

To run service locally: `make run`
Default port is 8080, but it can be altered in configs/default.yaml.

### Limitations

* Due to the technical requitements request timeout is implemented using middleware. However it seems to be not optimal for the issue, because timeout duration, on the on hand, is part of the buisness logic. On the other hand, chi and net/http server have built-in timeout support. 
* Go's standart log package is used for the means of simplicity, but if the app is meant to handle intence traffic, it's better to replace it with sturctured logging such as zap or logrus.  