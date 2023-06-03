package main

import (
	"context"
	"fmt"
	"log"
	"os"

	users "github.com/handson-go/chap8/user_service/service"
	"google.golang.org/grpc"
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

func main() {
	if len(os.Args) != 2 {
		log.Fatal(
			"Must specify a gRPC server address",
		)
	}
	conn, err := setupGrpcConn(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(
		c,
		&users.UserGetRequest{Email: "kim@naver.com"},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(
		os.Stdout, "User: %s %s \n",
		result.User.FirstName, result.User.LastName,
	)
}
