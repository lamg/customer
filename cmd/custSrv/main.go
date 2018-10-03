package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/lamg/customer"
	"google.golang.org/grpc"
)

type custSrv struct {
	customers []*pb.CustomerRequest
}

func (c *custSrv) CreateCustomer(ctx context.Context,
	in *pb.CustomerRequest) (r *pb.CustomerResponse, e error) {
	c.customers = append(c.customers, in)
	r = &pb.CustomerResponse{Id: in.Id, Success: true}
	return
}

func (c *custSrv) GetCustomers(filter *pb.CustomerFilter,
	stream pb.Customer_GetCustomersServer) (e error) {
	for i := 0; e == nil && i != len(c.customers); i++ {
		if strings.Contains(c.customers[i].Name, filter.Keyword) {
			e = stream.Send(c.customers[i])
		}
	}
	return
}

func main() {
	var addr string
	flag.StringVar(&addr, "p", ":50051",
		"Address the server listens incoming connections")
	flag.Parse()

	lst, e := net.Listen("tcp", addr)
	if e == nil {
		s := grpc.NewServer()
		pb.RegisterCustomerServer(s, new(custSrv))
		e = s.Serve(lst)
	}
	st := 0
	if e != nil {
		log.Fatalf("gRPC server: %s", e.Error())
		st = 1
	}
	os.Exit(st)
}
