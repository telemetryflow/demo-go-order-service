// Package handler contains command and query handlers.
package handler

import (
	"context"
	"fmt"
)

// Handler is a marker interface for all handlers
type Handler interface{}

// CommandBus dispatches commands to their handlers
type CommandBus struct {
	handlers map[string]interface{}
}

// NewCommandBus creates a new command bus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]interface{}),
	}
}

// Register registers a command handler
func (b *CommandBus) Register(commandType string, handler interface{}) {
	b.handlers[commandType] = handler
}

// Dispatch dispatches a command to its handler
func (b *CommandBus) Dispatch(ctx context.Context, commandType string, cmd interface{}) (interface{}, error) {
	handler, ok := b.handlers[commandType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for command type: %s", commandType)
	}

	// Handler must implement Handle(ctx, cmd) method
	h, ok := handler.(interface {
		Handle(context.Context, interface{}) (interface{}, error)
	})
	if !ok {
		return nil, fmt.Errorf("handler does not implement Handle method: %s", commandType)
	}

	return h.Handle(ctx, cmd)
}

// QueryBus dispatches queries to their handlers
type QueryBus struct {
	handlers map[string]interface{}
}

// NewQueryBus creates a new query bus
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]interface{}),
	}
}

// Register registers a query handler
func (b *QueryBus) Register(queryType string, handler interface{}) {
	b.handlers[queryType] = handler
}

// Dispatch dispatches a query to its handler
func (b *QueryBus) Dispatch(ctx context.Context, queryType string, q interface{}) (interface{}, error) {
	handler, ok := b.handlers[queryType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for query type: %s", queryType)
	}

	h, ok := handler.(interface {
		Handle(context.Context, interface{}) (interface{}, error)
	})
	if !ok {
		return nil, fmt.Errorf("handler does not implement Handle method: %s", queryType)
	}

	return h.Handle(ctx, q)
}
