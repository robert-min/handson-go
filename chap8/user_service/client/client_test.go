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

type dummyUserService struct {
	users.UnimplementedUsersServer
}

func (s *dummyUserService) GetUser(ctx context.Context, in *users.UserGetRequest) (*users.UserGetReply, error) {
	u := users.User{
		Id:        "user-123-a",
		FirstName: "kim",
		LastName:  "min",
		Age:       36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	users.RegisterUsersServer(s, &dummyUserService{})
	go func() {
		err := s.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return s, l
}

func TestGetUser(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	// 다이얼러 생성
	bufconnDialer := func(ctx context.Context, add string) (net.Conn, error) {
		return l.Dial()
	}

	// 테스트서버 연결을 위한 클라이언트 생성
	conn, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer), // 네트워크를 인메모리 연결로 구성
	)
	if err != nil {
		t.Fatal(err)
	}

	c := getUserServiceClient(conn)
	result, err := getUser(
		c,
		&users.UserGetRequest{Email: "kim@min.com"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.User.FirstName != "kim" || result.User.LastName != "min" {
		t.Fatalf(
			"Expected: kim doe, min: %s %s",
			result.User.FirstName,
			result.User.LastName,
		)
	}
}
