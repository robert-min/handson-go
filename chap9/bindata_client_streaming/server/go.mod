module github.com/handson-go/chap9/bindata_client_streaming/server

go 1.19

require (
	github.com/handson-go/chap9/bindata_client_streaming/service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.37.0
)

require (
	github.com/golang/protobuf v1.4.2 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/handson-go/chap9/bindata_client_streaming/service => ../service
