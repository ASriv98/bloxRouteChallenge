package main

import (
	"context"

	"sqsclientserver/src/data"
	"sqsclientserver/src/queue"
)

type Client interface {
	AddItem(data data.Data) error
	RemoveItem(key string) error
	GetItem(key string) error
	GetAllItems() error
}

const (
	_addItem     = "add-item"
	_deleteItem  = "delete-item"
	_getItem     = "get-item"
	_getAllItems = "get-all-items"
)

type client struct {
	q queue.Queue
}

func NewClient(q queue.Queue) client {
	return client{q: q}
}

func (c *client) AddItem(ctx context.Context, d data.Data) error {
	msg := data.NewMessage(
		_addItem,
		d)

	return c.q.SendMessage(ctx, msg)
}

func (c *client) DeleteItem(ctx context.Context, key string) error {
	msg := data.NewMessage(
		_deleteItem,
		data.NewData(key, nil))

	return c.q.SendMessage(ctx, msg)
}

func (c *client) GetItem(ctx context.Context, key string) error {
	msg := data.NewMessage(
		_getItem,
		data.NewData(key, nil))

	return c.q.SendMessage(ctx, msg)
}

func (c *client) GetAllItems(ctx context.Context) error {
	msg := data.NewMessage(
		_getAllItems,
		nil)

	return c.q.SendMessage(ctx, msg)
}
