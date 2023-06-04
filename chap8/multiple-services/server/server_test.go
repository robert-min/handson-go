package main

import (
	"context"
	"log"
	"net"
	"testing"

	svc "github.com/handson-go/chap8/multiple-services/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGRPCServer() (*grpc.Server, *bufconn.Listener) {
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
	s, l := startTestGRPCServer()
	defer s.GracefulStop()

	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	usersClient := svc.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&svc.UserGetRequest{
			Id:    "foo-bar",
			Email: "kim@do.com",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if resp.User.FirstName != "kim" {
		t.Errorf(
			"Expected FirstName to be: kim, Got: %s",
			resp.User.FirstName,
		)
	}
}

func TestRepoService(t *testing.T) {
	s, l := startTestGRPCServer()
	defer s.GracefulStop()

	bufconnDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	repoClient := svc.NewRepoClient(client)
	resp, err := repoClient.GetRepos(
		context.Background(),
		&svc.RepoGetRequest{
			CreateId: "user-11",
			Id:       "repo-11",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Repo) != 1 {
		t.Fatalf("Expected to get back 1 repo, got back: %d repos", len(resp.Repo))
	}
	gotId := resp.Repo[0].Id
	gotOwnerId := resp.Repo[0].Owner.Id

	if gotId != "repo-11" {
		t.Errorf(
			"Expected Repo ID to be: repo-11, Got: %s",
			gotId,
		)
	}

	if gotOwnerId != "user-11" {
		t.Errorf(
			"Expected Repo ID to be: user-11, Got: %s",
			gotOwnerId,
		)
	}
}
