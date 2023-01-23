package rsmq

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	e "github.com/galadd/rsmq-go/error"
)

const (
	q      = ":Q"
	queues = "QUEUES"
)

// Redis is the client of rsmq to execute queue and message operations
type Redis struct {
	client *redis.Client
	ns     string // Namespace
}

// Queue struct
type Queue struct {
	Name string
	Vt   int
	Delay int
	Maxsize int
	Created time.Time
	Modified time.Time
	TotalSent int 
	TotalReceived int
	Msgs int // Number of messages in the queue
	Hiddenmsgs int // Number of hidden messages in the queue
}

// Message struct
type Message struct {
	ID        string
	Message   string
	Receipt   string
	Visible   bool
	Received int // Number of times the message has been received
	Sent time.Time // Time when the message was sent
	FirstSeen time.Time
}

// NewRedis creates a new Redis client
func NewRedis(client *redis.Client, ns string) *Redis {
	return &Redis{
		client: client,
		ns:     ns,
	}
}

// CreateQueue creates a new queue with the given name
func (r *Redis) CreateQueue(queueName string, vt, delay uint, maxsize int) error {
	if queueName == "" {
		return e.ErrMissingParameter
	}
	
	if vt < 0 || vt > 9999999 {
		return e.ErrInvalidValue
	}
	if delay < 0 || delay > 9999999 {
		return e.ErrInvalidValue
	}
	if maxsize < 1024 || maxsize > 65536 {
		return e.ErrInvalidValue
	}

	// Check if queue already exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 1 {
		return e.ErrQueueExists
	}

	// Create queue
	r.client.HMSet(r.ns + ":" + queueName, map[string]interface{}{
		"vt": vt,
		"delay": delay,
		"maxsize": maxsize,
		"created": time.Now().Unix(),
		"modified": time.Now().Unix(),
		"totalsent": 0,
		"totalreceived": 0,
		"msgs": 0,
		"hiddenmsgs": 0,
	})
	r.client.SAdd(r.ns + ":" + queues, queueName)

	return nil
}

// DeleteQueue deletes the specified queue
func (r *Redis) DeleteQueue(queueName string) error {
	if queueName == "" {
		return e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return e.ErrQueueNotFound
	}

	// Delete queue
	r.client.Del(r.ns + ":" + queueName)
	r.client.SRem(r.ns + ":" + queues, queueName)

	return nil
}

// GetQueueAttributes returns the attributes of the specified queue
func (r *Redis) GetQueueAttributes(queueName string) (*Queue, error) {
	if queueName == "" {
		return nil, e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return nil, e.ErrQueueNotFound
	}

	// Get queue attributes
	q := &Queue{}
	q.Name = queueName
	q.Vt, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "vt").Val())
	q.Delay, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "delay").Val())
	q.Maxsize, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "maxsize").Val())
	q.Created, _ = time.Parse(r.client.HGet(r.ns + ":" + queueName, "created").Val(), "0")
	q.Modified, _ = time.Parse(r.client.HGet(r.ns + ":" + queueName, "modified").Val(), "0")
	q.TotalSent, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "totalsent").Val())
	q.TotalReceived, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "totalreceived").Val())
	q.Msgs, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "msgs").Val())
	q.Hiddenmsgs, _ = strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "hiddenmsgs").Val())

	return q, nil
}

// SetQueueAttributes sets the attributes of the specified queue
func (r *Redis) SetQueueAttributes(queueName string, vt, delay, maxsize int) error {
	if queueName == "" {
		return e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return e.ErrQueueNotFound
	}

	// Set queue attributes
	r.client.HSet(r.ns + ":" + queueName, "vt", vt)
	r.client.HSet(r.ns + ":" + queueName, "delay", delay)
	r.client.HSet(r.ns + ":" + queueName, "maxsize", maxsize)
	r.client.HSet(r.ns + ":" + queueName, "modified", time.Now().Unix())

	return nil
}

// GetQueueMessageCount returns the number of messages in the specified queue
func (r *Redis) GetQueueMessageCount(queueName string) (int, error) {
	if queueName == "" {
		return 0, e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return 0, e.ErrQueueNotFound
	}

	// Get queue message count
	count, _ := strconv.Atoi(r.client.HGet(r.ns + ":" + queueName, "msgs").Val())

	return count, nil
}

// ListQueues returns a list of all queues
func (r *Redis) ListQueues() ([]string, error) {
	return r.client.SMembers(r.ns + ":" + queues).Result()
}

// SendMessage sends a new message to the specified queue
func (r *Redis) SendMessage(queueName, message string) (string, error) {
	if queueName == "" || message == "" {
		return "", e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return "", e.ErrQueueNotFound
	}

	// Check if queue is full
	if r.client.HGet(r.ns + ":" + queueName, "msgs").Val() == r.client.HGet(r.ns + ":" + queueName, "maxsize").Val() {
		return "", e.ErrQueueFull
	}

	// Generate message ID
	id := uuid.NewV4().String()

	// Send message
	r.client.HMSet(r.ns + ":" + queueName + ":" + id, map[string]interface{}{
		"message": message,
		"sent": time.Now().Unix(),
		"firstseen": time.Now().Unix(),
	})
	r.client.HIncrBy(r.ns + ":" + queueName, "totalsent", 1)
	r.client.HIncrBy(r.ns + ":" + queueName, "msgs", 1)
	r.client.LPush(r.ns + ":" + queueName + ":Q", id)

	return id, nil
}

// ReceiveMessage receives a message from the specified queue
func (r *Redis) ReceiveMessage(queueName string) (*Message, error) {
	if queueName == "" {
		return nil, e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return nil, e.ErrQueueNotFound
	}

	// Check if queue is empty
	if r.client.HGet(r.ns + ":" + queueName, "msgs").Val() == "0" {
		return nil, e.ErrQueueEmpty
	}

	// Get message ID
	id, err := r.client.RPop(r.ns + ":" + queueName + ":Q").Result()
	if err != nil {
		return nil, err
	}

	// Get message
	m := &Message{}
	m.ID = id
	m.Message = r.client.HGet(r.ns + ":" + queueName + ":" + id, "message").Val()
	m.Sent, _ = time.Parse(r.client.HGet(r.ns + ":" + queueName + ":" + id, "sent").Val(), "0")
	m.FirstSeen, _ = time.Parse(r.client.HGet(r.ns + ":" + queueName + ":" + id, "firstseen").Val(), "0")

	// Set message attributes
	r.client.HSet(r.ns + ":" + queueName + ":" + id, "received", time.Now().Unix())
	r.client.HIncrBy(r.ns + ":" + queueName, "totalreceived", 1)
	r.client.HIncrBy(r.ns + ":" + queueName, "msgs", -1)
	r.client.HIncrBy(r.ns + ":" + queueName, "hiddenmsgs", 1)
	r.client.LPush(r.ns + ":" + queueName + ":R", id)

	return m, nil
}

// DeleteMessage deletes a message from the specified queue
func (r *Redis) DeleteMessage(queueName, id string) error {
	if queueName == "" || id == "" {
		return e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return e.ErrQueueNotFound
	}

	// Check if message exists
	if r.client.Exists(r.ns + ":" + queueName + ":" + id).Val() == 0 {
		return e.ErrMessageNotFound
	}

	// Delete message
	r.client.Del(r.ns + ":" + queueName + ":" + id)
	r.client.HIncrBy(r.ns + ":" + queueName, "msgs", -1)
	r.client.HIncrBy(r.ns + ":" + queueName, "hiddenmsgs", -1)
	r.client.LRem(r.ns + ":" + queueName + ":R", 0, id)

	return nil
}

// ChangeMessageVisibility changes the visibility of a message
func (r *Redis) ChangeMessageVisibility(queueName, id string, visibility int) error {
	if queueName == "" || id == "" {
		return e.ErrMissingParameter
	}

	// Check if queue exists
	if r.client.Exists(r.ns + ":" + queueName).Val() == 0 {
		return e.ErrQueueNotFound
	}

	// Check if message exists
	if r.client.Exists(r.ns + ":" + queueName + ":" + id).Val() == 0 {
		return e.ErrMessageNotFound
	}

	// Change message visibility
	r.client.HSet(r.ns + ":" + queueName + ":" + id, "visibility", visibility)
	r.client.HIncrBy(r.ns + ":" + queueName, "totalchanged", 1)

	return nil
}