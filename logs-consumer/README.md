# Client Logs Consumer
That is a simple client, that can consume logs from server based to choosen topics

You can pass all interested topics as command line arguments to the client

> example  
> will collect all messages from topic "error"
```shell
go run main.go error
```

One instance is launched at the start to write all logs to a file