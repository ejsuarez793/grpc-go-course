package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ejsuarez793/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a calculator")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	/*doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBiDiStreaming(c)*/

	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := calculatorpb.CalculatorRequest{Request: &calculatorpb.Calculator{FirstInteger: 1, SecondInteger: 2}}
	res, err := c.Sum(context.Background(), &req)

	if err != nil {
		log.Fatalf("There was an error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming RPC to decompose a number...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.Decompose(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Decomposition RPC %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}

		log.Printf("Response from Decomposition: %v", msg.GetPrimerNumber())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.Average(context.Background())

	requests := []*calculatorpb.AverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	if err != nil {
		log.Fatalf("error while calling Average %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("There was an error getting reponse from server %v", err)
	}

	fmt.Printf("Average Response: %v\n", res)

}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi stream RPC...")
	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("error while calling Average %v", err)
	}

	requests := []*calculatorpb.MaxRequest{
		{
			Number: 4,
		},
		{
			Number: 7,
		},
		{
			Number: 2,
		},
		{
			Number: 19,
		},
		{
			Number: 4,
		},
		{
			Number: 6,
		},
		{
			Number: 32,
		},
	}

	waitc := make(chan struct{})

	// we receive msg
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("There was an error receiving message", err)
				break
			}
			fmt.Printf("Received message: %v\n", res)
		}

		close(waitc)
	}()

	// we send msg
	go func() {

		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	// correct call
	doErrorCall(c, int32(10))

	// error call
	doErrorCall(c, int32(-2))

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	// correct call
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}

	}

	fmt.Printf("Result of square root of %v: %v\n", number, res.GetNumberRoot())
}
