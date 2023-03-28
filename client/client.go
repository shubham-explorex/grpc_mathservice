package main

import (
	"context"
	"fmt"
	"log"

	math1_v1 "example.com/proto"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	c := math1_v1.NewMathServiceClient(conn)

	num := math1_v1.Request{
		Num1: 200000,
		Num2: 0,
	}

	fmt.Println("before add")
	add, err := c.Add(context.Background(), &num)
	if err != nil {
		fmt.Println("err add")
		log.Fatalf("failed to call Add: %v", err)
	}
	fmt.Println("after add ")
	// fmt.Println("no error")
	fmt.Println(add.Result)

	sub, err := c.Subtract(context.Background(), &num)
	if err != nil {
		log.Fatalf("failed to call Subtract: %v", err)
	}
	fmt.Println(sub.Result)

	multiply, err := c.Multiply(context.Background(), &num)
	if err != nil {
		log.Fatalf("failed to call Multiply: %v", err)
	}
	fmt.Println(multiply.Result)

	div, err := c.Divide(context.Background(), &num)
	if err != nil {
		log.Fatalf("failed to call Divide: %v", err)
	}
	fmt.Println(div.Result)
}
