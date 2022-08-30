package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func newConn() (net.Conn, func()) {
	con, err := net.DialTimeout("tcp", "golang.org:http", time.Minute*3)
	if err != nil {
		log.Fatal("can't create client")
	}
	// log.Printf("local addr: %s, remote addr: %s\n", con.LocalAddr(), con.RemoteAddr())

	return con, func() {
		con.Close()
	}
}

func sendChunkedRequest(conn net.Conn, number int, timeout time.Duration) (string, error) {
	msg := []string{"POST / HTTP/1.1\r\n" +
		"Host: debian\r\n" +
		"Content-Type: text/html\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"4\r\n",
		"test77776767667\r\n",
		"0\r\n",
		"\r\n",
		"POST / HTTP/1.1\r\n" +
			"Host: debian\r\n" +
			"Content-Type: text/html\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"4\r\n",
		"test77776767667\r\n",
		"0\r\n",
		"\r\n",
	}
	for _, v := range msg {
		_, err := conn.Write([]byte(v))
		if err != nil {
			return "", fmt.Errorf("can't write, err %v", err)
		}
		// log.Printf("sent %v\n", v)

		time.Sleep(timeout)
	}

	reply := make([]byte, 1024)

	_, err := conn.Read(reply)
	if err != nil {
		return "", fmt.Errorf("can't read reply, err: %v", err)
	}

	replyStr := fmt.Sprintf("client %d got reply: %v", number, string(reply[:50]))
	fmt.Printf("%s\n", replyStr)
	return replyStr, nil
}

func main() {

	errCount := 0
	goodCount := 0
	wt := sync.WaitGroup{}

	for i := 1; i < 3000; i++ {
		wt.Add(1)

		go func(i int) {
			con, closeFunc := newConn()
			defer closeFunc()

			_, err := sendChunkedRequest(con, i, time.Second*10)
			if err != nil {
				errCount++
				log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
			} else {
				goodCount++
			}

			wt.Done()
		}(i)

		time.Sleep(time.Millisecond * 10)
	}

	wt.Wait()
	fmt.Printf("total errCount %v\ntotal goodCount %v\n", errCount, goodCount)
}
