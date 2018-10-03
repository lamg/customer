package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/lamg/customer"
	"google.golang.org/grpc"
)

func main() {
	var addr string
	flag.StringVar(&addr, "a", ":50051",
		"Address the client connects to")
	flag.Parse()

	c, e := grpc.Dial(addr, grpc.WithInsecure())
	defer c.Close()
	var cl pb.CustomerClient
	if e == nil {
		cl = pb.NewCustomerClient(c)
		e = create(cl)
	}
	if e == nil {
		e = get(cl, &pb.CustomerFilter{Keyword: ""})
	}

	code := 0
	if e != nil {
		code = 1
		log.Print(e.Error())
	}
	os.Exit(code)
}

func create(cl pb.CustomerClient) (e error) {
	for i := 0; e == nil && i != len(customers); i++ {
		var r *pb.CustomerResponse
		r, e = cl.CreateCustomer(context.Background(),
			customers[i])
		if e == nil {
			if !r.Success {
				e = fmt.Errorf("Not successful response")
			} else {
				log.Printf("Created new customer with ID: %d", r.Id)
			}
		} else {
			e = fmt.Errorf("failed to create customer %d: %s", i,
				e.Error())
		}
	}
	return
}

func get(cl pb.CustomerClient,
	filter *pb.CustomerFilter) (e error) {
	var stream pb.Customer_GetCustomersClient
	stream, e = cl.GetCustomers(context.Background(), filter)
	for e != io.EOF && e == nil {
		var cust *pb.CustomerRequest
		cust, e = stream.Recv()
		if e == nil {
			log.Printf("Customer %v", cust)
		}
	}
	if e == io.EOF {
		e = nil
	}
	return
}

var customers = []*pb.CustomerRequest{
	&pb.CustomerRequest{
		Id:    101,
		Name:  "Shiju Varghese",
		Email: "shiju@xyz.com",
		Phone: "732-757-2923",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "1 Mission Street",
				City:              "San Francisco",
				State:             "CA",
				Zip:               "94105",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "Greenfield",
				City:              "Kochi",
				State:             "KL",
				Zip:               "68356",
				IsShippingAddress: true,
			},
		},
	},
	&pb.CustomerRequest{
		Id:    102,
		Name:  "Irene Rose",
		Email: "irene@xyz.com",
		Phone: "732-757-2924",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "1 Mission Street",
				City:              "San Francisco",
				State:             "CA",
				Zip:               "94105",
				IsShippingAddress: true,
			},
		},
	},
}
