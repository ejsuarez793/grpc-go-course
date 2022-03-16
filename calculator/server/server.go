package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/ejsuarez793/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {

	fmt.Printf("Sum function was invoked with %v\n", req)

	a := req.GetRequest().FirstInteger
	b := req.GetRequest().SecondInteger

	res := &calculatorpb.CalculatorResponse{Result: a + b}

	return res, nil
}

func (*server) Decompose(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_DecomposeServer) error {
	fmt.Printf("Decompose function was invoked with %v\n", req)

	number := req.GetNumber()

	k := int32(2)

	for {
		if number <= 1 {
			break
		}

		if number%k == 0 {
			fmt.Println("One factor is: ", k)
			res := &calculatorpb.PrimeNumberDecompositionReponse{PrimerNumber: k}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			number = number / k
		} else {
			k = k + 1
		}

	}

	return nil

}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Println("Average function was invoked with a streaming request")
	s := float32(0)
	n := float32(0)
	number := float32(0)
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// we have finished reading the client stream
			avg := s / n
			res := &calculatorpb.AverageResponse{Avg: float32(avg)}
			return stream.SendAndClose(res)

		}

		if err != nil {
			log.Fatalf("there was an error reading client stream %v", err)
		}

		number = float32(req.GetNumber())

		fmt.Println("received number: ", number)

		s += number
		n++
	}
}

func (*server) Max(stream calculatorpb.CalculatorService_MaxServer) error {

	fmt.Println("Max method was invoked with a streaming request")

	current_maximun := int32(math.MinInt32)
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("There was an error receiving client stream message", err)
			return err
		}

		number := req.GetNumber()

		fmt.Printf("Received number: %v\n", number)

		if number > current_maximun {
			current_maximun = number
			sendErr := stream.Send(&calculatorpb.MaxResponse{Number: current_maximun})
			if sendErr != nil {
				log.Fatalf("There was an error sending data to client: %v", sendErr)
				return sendErr
			}
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Received SquareRoot RPC: %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(number))}, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	// open connection in port
	// handle error
	//register grpc service

	//handle request

}
