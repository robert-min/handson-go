package main

import (
	"context"
	"log"
	"net"
	"testing"

	users "github.com/handson-go/chap8/user_service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Start test GRPC server with real network listener
func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestUserService(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	// 다이얼러 생성
	bufconnDialer := func(ctx context.Context, add string) (net.Conn, error) {
		return l.Dial()
	}

	// 테스트서버 연결을 위한 클라이언트 생성
	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer), // 네트워크를 인메모리 연결로 구성
	)
	if err != nil {
		t.Fatal(err)
	}

	usersClient := users.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&users.UserGetRequest{
			Email: "kim@naver.com",
			Id:    "foo-bar",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "kim" {
		t.Errorf(
			"Expected FirstName to be : kim, Got : %s",
			resp.User.FirstName,
		)
	}
}
