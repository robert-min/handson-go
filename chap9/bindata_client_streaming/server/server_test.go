package main

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"

	svc "github.com/handson-go/chap9/bindata_client_streaming/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() *bufconn.Listener {
	l := bufconn.Listen(1)
	s := grpc.NewServer()
	registerServices(s)
	go func() {
		startServer(s, l)
	}()
	return l
}

func TestCreateRepo(t *testing.T) {
	l := startTestGrpcServer()

	bufconnDialer := func(
		ctx context.Context, addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"", grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal("DialContext", err)
	}

	// 1. 리포지터리 이름과 소유자 정보를 주고 받음
	repoClient := svc.NewRepoClient(client)
	stream, err := repoClient.CreateRepo(
		context.Background(),
	)
	if err != nil {
		t.Fatal("CreateRepo", err)
	}

	c := svc.RepoCreateRequest_Context{
		Context: &svc.RepoContext{
			CreateId: "user-123",
			Name:     "test-repo",
		},
	}
	r := svc.RepoCreateRequest{
		Body: &c,
	}
	err = stream.Send(&r)
	if err != nil {
		t.Fatal("StreamSend", err)
	}

	// 2. data 필드만 설정한 RepoCreateContext 객체를 전송
	data := "Arbitray data bytes"
	repoData := strings.NewReader(data)

	for {
		b, err := repoData.ReadByte()
		if err == io.EOF {
			break
		}

		bData := svc.RepoCreateRequest_Data{
			Data: []byte{b},
		}
		r := svc.RepoCreateRequest{
			Body: &bData,
		}
		err = stream.Send(&r)
		if err != nil {
			t.Fatal("StreaSend", err)
		}
		l.Close()
	}

	// 3. 서버에서 받아들인 응답이 정상인지 확인
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal("CloseAndRecv", err)
	}
	expectedSize := int32(len(data))
	if resp.Size != expectedSize {
		t.Errorf(
			"Expected Repo Created to be: %d bytes Got back: %d",
			expectedSize,
			resp.Size,
		)
	}
	expectedRepoUrl := "https://example.com/user-123/test-repo"
	if resp.Repo.Url != expectedRepoUrl {
		t.Errorf(
			"Expected Repo URL to be: %s, Got: %s",
			expectedRepoUrl,
			resp.Repo.Url,
		)
	}
}
