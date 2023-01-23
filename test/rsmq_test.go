package test

import (
	"testing"
	"fmt"

	"github.com/galadd/rsmq-go"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var client *redis.Client
var r *rsmq.Redis

func init() {
	s, err := miniredis.Run()
    if err != nil {
        panic(err)
    }

    client = redis.NewClient(&redis.Options{Addr: s.Addr()})
	r = rsmq.NewRedis(client, "rsmq")
}

// TestCreateQueue tests the creation of a queue
func TestCreateQueue(t *testing.T) {
	// Create a queue
	err := r.CreateQueue("test", 30, 0, 65536)
	assert.Nil(t, err)
	// print checkmark and success message if no error
	fmt.Println("\u2713", "CreateQueue")

	// Create a queue with an invalid name
	err = r.CreateQueue("", 30, 0, 65536)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when creating a queue with an invalid name")

	// Create a queue with an invalid vt
	err = r.CreateQueue("test", 99999999, 0, 65536)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when creating a queue with an invalid vt")

	// Create a queue with an invalid delay
	err = r.CreateQueue("test", 30, 99999999, 65536)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when creating a queue with an invalid delay")

	// Create a queue with an invalid maxsize
	err = r.CreateQueue("test", 30, 0, 65537)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when creating a queue with an invalid maxsize")
}

// TestSetQueueAttributes tests the setting of queue attributes
func TestSetQueueAttributes(t *testing.T) {
	// Set queue attributes
	err := r.SetQueueAttributes("test", 30, 0, 65535)
	assert.Nil(t, err)
	fmt.Println("\u2713", "SetQueueAttributes")

	// Set Queue atrributes of a non-existing queue
	err = r.SetQueueAttributes("test2", 30, 0, 65535)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when setting queue attributes of a non-existing queue")
}

// TestGetQueueAttributes tests the getting of queue attributes
func TestGetQueueAttributes(t *testing.T) {
	// Get queue attributes
	q, err := r.GetQueueAttributes("test")
	assert.Nil(t, err)
	fmt.Println("\u2713", "GetQueueAttributes")
	assert.Equal(t, 30, q.Vt)
	assert.Equal(t, 0, q.Delay)
	assert.Equal(t, 65535, q.Maxsize)
	fmt.Println("\u2713", "Returns the correct queue attributes")

	// Get queue attributes of a non-existing queue
	q, err = r.GetQueueAttributes("test2")
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when getting queue attributes of a non-existing queue")
}

// TestListQueues tests the listing of queues
func TestListQueues(t *testing.T) {
	// List queues
	queues, err := r.ListQueues()
	assert.Nil(t, err)
	fmt.Println("\u2713", "ListQueues")
	assert.Equal(t, 1, len(queues))
	assert.Equal(t, "test", queues[0])
	fmt.Println("\u2713", "Returns the correct queues")
}

// TestSendMessage tests the sending of a message to a queue
func TestSendMessage(t *testing.T) {
	// Send a message
	id, err := r.SendMessage("test", "test message")
	assert.Nil(t, err)
	fmt.Println("\u2713", "SendMessage")
	assert.Equal(t, 36, len(id))
	fmt.Println("\u2713", "Returns a valid id")

	// Send a message to a non-existing queue
	id, err = r.SendMessage("test2", "test message")
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when sending a message to a non-existing queue")
}

// TestReceiveMessage tests the receiving of a message from a queue
func TestReceiveMessage(t *testing.T) {
	// Receive a message
	m, err := r.ReceiveMessage("test")
	assert.Nil(t, err)
	fmt.Println("\u2713", "ReceiveMessage")
	assert.Equal(t, "test message", m.Message) //
	fmt.Println("\u2713", "Returns the correct message")

	// Receive a message from a non-existing queue
	m, err = r.ReceiveMessage("test2")
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when receiving a message from a non-existing queue")
}

// TestDeleteMessage tests the deletion of a message from a queue
func TestDeleteMessage(t *testing.T) {
	// send a message
	id, _ := r.SendMessage("test", "test message")

	// Delete a message
	err := r.DeleteMessage("test", id)
	assert.Nil(t, err)
	fmt.Println("\u2713", "DeleteMessage")

	// Delete a message from a non-existing queue
	err = r.DeleteMessage("test2", id)
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when deleting a message from a non-existing queue")
}

// TestDeleteQueue tests the deletion of a queue
func TestDeleteQueue(t *testing.T) {
	// Delete a queue
	err := r.DeleteQueue("test")
	assert.Nil(t, err)
	fmt.Println("\u2713", "DeleteQueue")

	// Delete a non-existing queue
	err = r.DeleteQueue("test2")
	assert.NotNil(t, err)
	fmt.Println("\u2713", "Returns an error when deleting a non-existing queue")
}