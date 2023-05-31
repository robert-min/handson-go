package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

type requestContextKey struct{} // key경우 string과 같은 기본 데이터 타입이 안되기 때문에 공백의 구조체로 정의
type requestContextValue struct {
	requestID string
}

// Add Request ID
func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestID: requestID,
	}
	currentCtx := r.Context()
	// WithValue(저장될 값의 콘텍스트를 식별하기 위한 부모 Context 객체, 맵 자료구조에서 데이터 키값을 식별하기 위한 interface 객체,
	// 데이터 자체의 interface{} 객체)
	newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
	return r.WithContext(newCtx)
}

func logRequest(r *http.Request) {
	ctx := r.Context()
	v := ctx.Value(requestContextKey{})

	if m, ok := v.(requestContextValue); ok {
		log.Printf("Processing request: %s", m.requestID)
	}
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprintf(w, "Request processed")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	requestID := "request-123-abc"
	r = addRequestID(r, requestID)
	processRequest(w, r)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandler)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
