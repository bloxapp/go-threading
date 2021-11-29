package queue

import (
	"github.com/bloxapp/go-threading/channel"
	"github.com/bloxapp/go-threading/queue/policies"
)

type ItemState int

const (
	ItemPopped    ItemState = 1
	ItemCancelled ItemState = 2
)

type Item interface {
	statefullItem
	PolicyManager() policies.PolicyManager
	Item() interface{}
	// Waiter will fire if the item was popped or cancelled
	Waiter() *channel.Waiter
}

type statefullItem interface {
	Popped()
	Cancelled()
}

type item struct {
	item    interface{}
	waiter  *channel.Waiter
	manager policies.PolicyManager
}

func NewItem(i interface{}, policyManager policies.PolicyManager) Item {
	return &item{
		item:    i,
		manager: policyManager,
		waiter:  channel.NewWaiter(),
	}
}

func (i *item) PolicyManager() policies.PolicyManager {
	return i.manager
}

func (i *item) Item() interface{} {
	return i.item
}

func (i *item) Waiter() *channel.Waiter {
	return i.waiter
}

func (i *item) Popped() {
	i.waiter.Fire(ItemPopped)
}

func (i *item) Cancelled() {
	i.waiter.Fire(ItemCancelled)
}
