package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// longRunningProccess is example to send data
func longRunningProcess(logWriter *io.PipeWriter) {
	for i := 0; i <= 20; i++ {
		fmt.Fprintf(logWriter, `{"id": %d, "user_ip": "172.121.19.21", "event": "click_on_add_cart" }`, i)
		fmt.Fprintln(logWriter)
		time.Sleep(1 * time.Second)
	}
}

// progresStreamer process response
func progressStreamer(logReader *io.PipeReader, w http.ResponseWriter, done chan struct{}) {
	buf := make([]byte, 500) // 특정 시점에 파이프에서 읽어 들일 최대 크기를 의미

	// 응답 데이터가 곧바고 클라이언트에게 Flush 매서드 호출 전 Flusher 인터페이스 구현 확인
	// f : http.flusher 객체, flushSupported : true
	f, flushSupported := w.(http.Flusher)

	defer logReader.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff") // 클라이언트 측에 데이터를 버퍼링하지 않도록 브라우저에 알림

	for {
		n, err := logReader.Read(buf)
		if err == io.EOF {
			break
		}
		w.Write(buf[:n])
		if flushSupported {
			f.Flush()
		}
	}
	done <- struct{}{}
}

func longRunningProcessHandler(w http.ResponseWriter, r *http.Request) {

	done := make(chan struct{})
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter) // 오래 걸리는 잡 실행
	go progressStreamer(logReader, w, done)

	<-done
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/job", longRunningProcessHandler)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("Server could not start listening on %s. Error: %v", listenAddr, err)
	}
}
