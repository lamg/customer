# Customers: a gRPC example

Example gRPC client and server written in Go, following https://medium.com/@shijuvar/building-high-performance-apis-in-go-using-grpc-and-protocol-buffers-2eda5b80771b

The file `customer.pb.go` is generated by command `protoc -I . customer.proto --go_out=plugins=grpc:.` executed in the project's root directory.

## Install and usage

```sh
git clone https://github.com/lamg/customer
cd customer/cmd/custSrv && go install && cd ../../..
cd customer/cmd/custCl && go install
custSrv&
custCl
pkill custSrv
```