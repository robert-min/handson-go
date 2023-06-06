package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	svc "github.com/handson-go/chap9/interceptors/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func metadataUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctxWithMetatdata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123",
	)
	return invoker(ctxWithMetatdata, method, req, reply, cc, opts...)
}

func metadataStreamInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123",
	)
	clientStream, err := streamer(
		ctxWithMetadata,
		desc,
		cc,
		method,
		opts...,
	)
	return clientStream, err
}

// Set up Grpc Connection with Server
func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(), // 공백의 Context
		addr,                 // 연결할 서버의 주소
		grpc.WithInsecure(),  // 전송계층보안 사용하지 않고 서버와 통신을 수립하기 위해 사용
		grpc.WithBlock(),     // 함수가 반환되기 전 연결이 먼저 수립외더 연결 객체가 반환되도록 설정
		grpc.WithUnaryInterceptor(metadataUnaryInterceptor),   // 단항 클라이언트 인터셉트
		grpc.WithStreamInterceptor(metadataStreamInterceptor), // 스트림 클라이언트 인터셉트
	)
}

func getUserServiceClient(conn *grpc.ClientConn) svc.UsersClient {
	return svc.NewUsersClient(conn)
}

func getUser(
	client svc.UsersClient,
	u *svc.UserGetRequest,
) (*svc.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func setupChat(r io.Reader, w io.Writer, c svc.UsersClient) error {
	stream, err := c.GetHelp(context.Background())
	if err != nil {
		return err
	}

	for {
		scanner := bufio.NewScanner(r)
		prompt := "Request: "
		fmt.Fprint(w, prompt)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}

		msg := scanner.Text()
		if msg == "quit" {
			break
		}
		request := svc.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&request)
		if err != nil {
			return err
		}

		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("Response: %s\n", resp.Response)
	}
	return stream.CloseSend()
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal(
			"Must specify a gRPC server address",
		)
	}
	serverAddr := os.Args[1]
	methodName := os.Args[2]

	conn, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		result, err := getUser(
			c,
			&svc.UserGetRequest{Email: "kim@naver.com"},
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(
			os.Stdout, "User: %s %s \n",
			result.User.FirstName, result.User.LastName,
		)
	case "GetHelp":
		err = setupChat(os.Stdin, os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}

}
