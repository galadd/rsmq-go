# RSMQ (Redis Simple Message Queue) implementation in Go

This is an implementation of the RSMQ (Redis Simple Message Queue) in Go. RSMQ is a lightweight message queue that uses Redis as a backend.

## Features
- Simple to use
- High performance
- Lightweight
- Redis backend

## Installation
To use this RSMQ implementation in your Go project, you need to have Go and Redis installed. You can install the package by running the following command:

```bash
go get github.com/galadd/rsmq-go
```

You can also add the package as a dependency in your go.mod file:

```go
require github.com/galadd/rsmq-go v1.0.0
```

## Usage
Here is a simple example of how to use the package:

```go
package main

import (
    "fmt"
    "time"

    "github.com/galadd/rsmq-go"
)

func main() {
    // Create a new redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

   // Create a new Redis struct with the client and namespace
	r := NewRedis(client, "rsmq")

    // Create a new queue with the name "testqueue", 
    // visibility timeout of 30 seconds, delay of 0 seconds, and maximum size of 4096 bytes
	err := r.CreateQueue("testqueue", 30, 0, 4096)
	if err != nil {
		fmt.Println(err)
	} 
	fmt.Println("Queue 'testqueue' created successfully.")

    // Get Queue Attributes
    attrs, err := r.GetQueueAttributes("testqueue")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Queue Attributes:", attrs)

    // Send a message to the queue
    msgID, err := r.SendMessage("testqueue", "Hello World!")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Message sent with ID:", msgID)

    // Receive a message from the queue
    msg, err := r.ReceiveMessage("testqueue")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Received message:", msg.Message)

    // Delete the Queue
    err = r.DeleteQueue("testqueue")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Queue 'testqueue' deleted successfully.")
}
```

The package provides several methods such as CreateQueue, SendMessage, ReceiveMessage, ChangeMessageVisibility, DeleteMessage, DeleteQueue and more. You can find more information about the methods in the package documentation.

## Contributing
We welcome contributions to this RSMQ implementation in Go. If you would like to contribute, please fork the repository and submit a pull request.

## License
This RSMQ implementation in Go is released under the MIT License.

## Contact
If you have any questions or feedback, please feel free to open an issue on the repository.

## Acknowledgements
This implementation is based on the original RSMQ (Redis Simple Message Queue) project. More information about the original project can be found at [smrchy/rsmq](https://github.com/smrchy/rsmq.git).

## Conclusion
With this RSMQ (Redis Simple Message Queue) implementation in Go, you can easily add a message queue to your Go projects and take advantage of the high performance and scalability of Redis. The package provides an easy-to-use interface for interacting with Redis and is suitable for a wide range of use cases.