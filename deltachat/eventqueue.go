package deltachat

import (
	"context"
	"sync"
)

type eventQueue struct {
	sync.Mutex
	items      []Event
	changedCtx context.Context
	changed    context.CancelFunc
	closedCtx  context.Context
	close      context.CancelFunc
}

func newEventQueue() *eventQueue {
	closedCtx, close := context.WithCancel(context.Background())
	changedCtx, changed := context.WithCancel(context.Background())
	return &eventQueue{changedCtx: changedCtx, changed: changed, closedCtx: closedCtx, close: close}
}

func (self *eventQueue) Put(item Event) {
	self.Lock()
	defer self.Unlock()

	select {
	case <-self.closedCtx.Done():
	default:
		self.items = append(self.items, item)
		self.changed()
	}
}

func (self *eventQueue) Close() {
	self.Lock()
	defer self.Unlock()

	select {
	case <-self.closedCtx.Done():
	default:
		self.close()
		self.items = nil
	}
}

func (self *eventQueue) Pop(ctx context.Context) Event {
	select {
	case <-ctx.Done():
		return nil
	case <-self.closedCtx.Done():
		return nil
	case <-self.changedCtx.Done():
		self.Lock()
		defer self.Unlock()

		select {
		case <-self.closedCtx.Done():
			return nil
		default:
		}

		item := self.items[0]
		self.items[0] = nil
		self.items = self.items[1:]

		if len(self.items) == 0 {
			self.changedCtx, self.changed = context.WithCancel(context.Background())
		}

		return item
	}
}
