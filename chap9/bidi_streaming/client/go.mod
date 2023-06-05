module github.com/handson-go/chap9/bidi_streaming/client

go 1.19

replace github.com/handson-go/chap9/bidi_streaming/service => ../service

require (
	github.com/handson-go/chap9/bidi_streaming/service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.37.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
