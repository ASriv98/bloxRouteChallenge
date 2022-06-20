package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"sqsclientserver/config"
	"sqsclientserver/src/data"
	"sqsclientserver/src/logging"
	"sqsclientserver/src/queue"
)

const (
	_addItem     = "add-item"
	_deleteItem  = "delete-item"
	_getItem     = "get-item"
	_getAllItems = "get-all-items"

	_maxWorkers = 4
	_idleTime   = 500
)

type Server interface {
	Start() error
	Stop() error
}

type server struct {
	items        sync.Map
	idleInterval time.Duration
	numWorkers   int

	q queue.Queue
	l *logging.Logger
}

func NewServer(queue queue.Queue, l *logging.Logger, d config.Data) server {
	var (
		idleInterval time.Duration
		idleTime     int64
		numWorkers   int64
	)

	if d.ServerIdleInterval != "" {
		idleTime, _ = strconv.ParseInt(d.ServerIdleInterval, 10, 64)
		idleInterval = time.Duration(idleTime) * time.Millisecond

	} else {
		idleTime := _idleTime
		idleInterval = time.Duration(idleTime) * time.Millisecond
	}

	if d.NumServerWorkers != "" {
		numWorkers, _ = strconv.ParseInt(d.NumServerWorkers, 10, 64)
	} else {
		numWorkers = _maxWorkers
	}

	return server{
		items:        sync.Map{},
		q:            queue,
		idleInterval: idleInterval,
		numWorkers:   int(numWorkers),
		l:            l,
	}
}

func (s *server) Start(ctx context.Context) {
	s.l.Debugf("server starting with %d workers...", s.numWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(s.numWorkers)

	for id := 1; id <= s.numWorkers; id++ {
		go s.worker(ctx, wg, id)
	}

	wg.Wait()

	s.l.Debugf("server closing...")
}

func (s *server) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		recv, handle, err := s.q.ReceiveMessage(ctx)
		if err != nil {
			s.l.Errorf("worker #%d errored : %w", id, err)
			time.Sleep(s.idleInterval)

			continue
		} else {
			// here we will parse the message that is returned from the queue
			if recv != nil {
				s.l.Debugf("\nreceived queue message")
				var msg data.Message

				bytes := []byte(*recv.(*string))

				if err = json.Unmarshal(bytes, &msg); err != nil {
					s.l.Errorf("\n error while receiving message  %v", err)
					s.l.Errorf("unmarshal : %v", err)

					err = s.q.DeleteMessage(ctx, handle)
					if err != nil {
						s.l.Errorf("delete queue message %v", err)
					}
				}

				// each received message gets processed in a separate goroutine
				go func() {
					s.l.Debugf("processing message request to %v", msg.Method)
					s.l.Debugf("message id is %v", msg.ID)

					err := s.processMessage(msg)
					if err != nil {
						s.l.Errorf("process queue message with id %v %v", msg.ID, err)
					}

					s.l.Debugf("deleting message in queue post processing with id %v", msg.ID)

					defer func(q queue.Queue, ctx context.Context, handle *string) {
						err := q.DeleteMessage(ctx, handle)
						if err != nil {

						}
					}(s.q, ctx, handle)
				}()
			}
		}
	}
}

func (s *server) processMessage(msg data.Message) error {
	switch msg.Method {

	case _addItem:
		key, ok := msg.Params.(map[string]interface{})["Key"].(string)
		value, ok2 := msg.Params.(map[string]interface{})["Value"]
		if !ok || !ok2 {
			return fmt.Errorf("problems with serialized item")
		}
		item := data.NewData(key, value)
		s.addItem(&item)

	case _deleteItem:
		key, ok := msg.Params.(map[string]interface{})["Key"].(string)
		if !ok {
			return fmt.Errorf("problems with serialized item key")
		}
		s.deleteItem(key)

	case _getItem:
		fmt.Printf("%+v msg", msg.Params)
		key, ok := msg.Params.(map[string]interface{})["Key"]
		if !ok {
			return fmt.Errorf("problems with serialized item key")
		}

		s.getItem(key.(string))

	case _getAllItems:
		s.getAllItems()

	default:
		return fmt.Errorf("wrong method call in message")
	}

	return nil
}

func (s *server) addItem(data *data.Data) {
	s.l.Infof("adding key %v with value %v to server", data.Key, data.Value)
	s.items.Store(data.Key, data.Value)
	str := fmt.Sprintf("add %s = %s", data.Key, data.Value)
	s.l.Infof(str)
}

func (s *server) deleteItem(key string) {
	s.l.Infof("deleting item with key %v from server", key)
	s.items.Delete(key)
	str := fmt.Sprintf("delete %s", key)
	s.l.Infof(str)
}

func (s *server) getItem(key string) {
	s.l.Debugf("getting single item with key %v from server", key)
	item, ok := s.items.Load(key)
	if !ok {
		str := fmt.Sprintf("no item exists for key %s", key)
		s.l.Infof(str)
	} else {
		str := fmt.Sprintf("get %s = %s", key, item)
		s.l.Infof(str)
	}
}

func (s *server) getAllItems() {
	s.l.Infof("getting all items from server")
	var data string
	var count int

	s.items.Range(
		func(key interface{}, value interface{}) bool {
			data += fmt.Sprintf("{%s: %s} ", key, value)
			count++
			return true
		})

	print := fmt.Sprintf("\ntotal %d items :", count, data)
	fmt.Println(print)
}
