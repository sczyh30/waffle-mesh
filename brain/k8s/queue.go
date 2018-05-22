package k8s

import (
	"sync"
	"time"
	"k8s.io/client-go/util/flowcontrol"
	"log"
)

// Here we define the queue because queue.Interface cannot support EventTask.
type RateLimitingWorkingQueue struct {
	delay   time.Duration
	queue   []EventTask
	mutex    sync.Mutex
	closing bool
}

func NewQueue(delay time.Duration) *RateLimitingWorkingQueue {
	return &RateLimitingWorkingQueue{
		delay:   delay,
		queue:   make([]EventTask, 0),
		closing: false,
		mutex:    sync.Mutex{},
	}
}

func (q *RateLimitingWorkingQueue) Push(task EventTask) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if !q.closing {
		q.queue = append(q.queue, task)
	}
}

func (q *RateLimitingWorkingQueue) Run(stop <-chan struct{}) {
	go func() {
		<-stop
		q.mutex.Lock()
		q.closing = true
		q.mutex.Unlock()
	}()

	rateLimit := 100
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(float32(rateLimit), 10*rateLimit)

	var task EventTask
	for {
		rateLimiter.Accept()

		q.mutex.Lock()
		if q.closing {
			q.mutex.Unlock()
			return
		} else if len(q.queue) == 0 {
			q.mutex.Unlock()
		} else {
			task, q.queue = q.queue[0], q.queue[1:]
			q.mutex.Unlock()

			for {
				err := task.f(task.obj, task.key, task.event)
				if err != nil {
					log.Printf("Work task failed (%v), repeating after delay %v\n", err, q.delay)
					time.Sleep(q.delay)
				} else {
					break
				}
			}
		}
	}
}

