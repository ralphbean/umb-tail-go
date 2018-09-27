Try it out:

```bash
$ go get github.com/go-stomp/stomp
$ go build umb-tail.go
$ ./umb-tail \
    -cert /path/to/umb-cert.pem \
    -key  /path/to/umb-key.pem \
    -queue "/queue/Consumer.USERNAME.go-test.VirtualTopic.eng.>"