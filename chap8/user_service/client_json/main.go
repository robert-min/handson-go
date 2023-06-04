package main

import (
	"context"
	"fmt"
	"log"
	"os"

	users "github.com/handson-go/chap8/user_service/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Set up Grpc Connection with Server
func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(), // 공백의 Context
		addr,                 // 연결할 서버의 주소
		grpc.WithInsecure(),  // 전송계층보안 사용하지 않고 서버와 통신을 수립하기 위해 사용
		grpc.WithBlock(),     // 함수가 반환되기 전 연결이 먼저 수립외더 연결 객체가 반환되도록 설정
	)
}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func getUser(
	client users.UsersClient,
	u *users.UserGetRequest,
) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func createUserRequest(jsonQuery string) (*users.UserGetRequest, error) {
	u := users.UserGetRequest{}
	input := []byte(jsonQuery)
	return &u, protojson.Unmarshal(input, &u)
}

func getUserResponseJson(result *users.UserGetReply) ([]byte, error) {
	return protojson.Marshal(result)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal(
			"Must specify a gRPC server address and search query",
		)
	}
	serverAddr := os.Args[1]

	u, err := createUserRequest(os.Args[2])
	if err != nil {
		log.Fatalf("Bad user input : %v", err)
	}

	conn, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(
		c,
		u,
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := getUserResponseJson(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(
		os.Stdout, string(data),
	)
}
