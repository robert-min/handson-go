package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	svc "github.com/handson-go/chap8/multiple-services/service"
	"google.golang.org/grpc"
)

type userSerivce struct {
	svc.UnimplementedUsersServer
}

type repoService struct {
	svc.UnimplementedRepoServer
}

func (s *userSerivce) GetUser(ctx context.Context, in *svc.UserGetRequest) (*svc.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	componets := strings.Split(in.Email, "@")
	if len(componets) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := svc.User{
		Id:        in.Id,
		FirstName: componets[0],
		LastName:  componets[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: &u}, nil
}

func (s *repoService) GetRepos(ctx context.Context, in *svc.RepoGetRequest) (*svc.RepoGetReply, error) {
	log.Printf(
		"Received request for repo with CreateId: %s Id: %s\n",
		in.CreateId,
		in.Id,
	)
	repo := svc.Repository{
		Id:    in.Id,
		Name:  "test repo",
		Url:   "https://git.example.com/test/repo",
		Owner: &svc.User{Id: in.CreateId, FirstName: "kim"},
	}
	r := svc.RepoGetReply{
		Repo: []*svc.Repository{&repo},
	}
	return &r, nil
}

func registerServices(s *grpc.Server) {
	svc.RegisterUsersServer(s, &userSerivce{})
	svc.RegisterRepoServer(s, &repoService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
