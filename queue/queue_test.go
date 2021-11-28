package queue

import (
	"fmt"
	policies2 "go-threading/queue/policies"
	"go-threading/threadsafe"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/require"
)

var evictImmediatelyF = func(i ...interface{}) func() policies2.Policy {
	return func() policies2.Policy {
		return newEvictImmediately()
	}
}

type evictImmediately struct {
}

func newEvictImmediately() policies2.Policy {
	return &evictImmediately{}
}

func (tp *evictImmediately) Evacuate() bool {
	return true
}

func (tp *evictImmediately) Item() interface{} {
	return nil
}

func TestQueueAddAndPop(t *testing.T) {
	q := New(FIFO, 10)
	require.True(t, q.Add(true, ""))
	require.True(t, q.Pop(DefaultItemIndex).(bool))
	require.EqualValues(t, 0, q.Len())
}

func TestQueueIndex(t *testing.T) {
	t.Run("pop from non existing index", func(t *testing.T) {
		q := New(FIFO, 10)
		require.Nil(t, q.Pop("non_existing_index"))
	})

	t.Run("default index", func(t *testing.T) {
		q := New(FIFO, 10)
		q.Add("item", "")
		require.EqualValues(t, "item", q.Pop(DefaultItemIndex))
	})

	t.Run("add new index", func(t *testing.T) {
		q := New(FIFO, 10)
		q.Add("item", "index")
		require.Nil(t, q.Pop(DefaultItemIndex))
		require.EqualValues(t, "item", q.Pop("index"))
		require.Nil(t, q.(*queue).queue["index"])
	})

}

func TestEmptyQueue(t *testing.T) {
	q := New(FIFO, 10)
	require.Nil(t, q.Pop(DefaultItemIndex))
}

func TestQueueDirection(t *testing.T) {
	t.Run("fifo", func(t *testing.T) {
		q := New(FIFO, 10)

		q.Add("first", "")
		q.Add("second", "")
		q.Add("third", "")

		require.EqualValues(t, "first", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "second", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "third", q.Pop(DefaultItemIndex))
	})

	t.Run("lifo", func(t *testing.T) {
		q := New(LIFO, 10)

		q.Add("first", "")
		q.Add("second", "")
		q.Add("third", "")

		require.EqualValues(t, "third", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "second", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "first", q.Pop(DefaultItemIndex))
	})
}

func TestPopWait(t *testing.T) {
	q := New(FIFO, 10)

	go func() {
		time.Sleep(time.Millisecond * 60)
		q.Add(true, "")
		time.Sleep(time.Millisecond * 60)
		q.Add(1, "")
		time.Sleep(time.Millisecond * 60)
		q.Add("t", "")
	}()

	require.True(t, q.PopWait(DefaultItemIndex).Wait().(bool))
	require.EqualValues(t, 1, q.PopWait(DefaultItemIndex).Wait().(int))
	require.EqualValues(t, "t", q.PopWait(DefaultItemIndex).Wait().(string))
}

func TestAddWhenFull(t *testing.T) {
	t.Run("multiple adds > capacity", func(t *testing.T) {
		q := New(FIFO, 3)
		require.True(t, q.Add("item", "index"))
		require.True(t, q.Add("item", "index2"))
		require.True(t, q.Add("item", "index3"))
		require.False(t, q.Add("item", "index4"))
	})
	t.Run("multiple indexes on immediate eviction", func(t *testing.T) {
		q := New(FIFO, 3, evictImmediatelyF())
		for i := 0; i < 100; i++ {
			require.True(t, q.Add("item2", Index(fmt.Sprintf("index_%d", i))))
		}
	})
}

func TestAddStateful(t *testing.T) {
	t.Run("fired when popped", func(t *testing.T) {
		q := New(FIFO, 3)
		res, waiter := q.AddStateful("item", "index")
		require.True(t, res)

		called := threadsafe.Int32(0)
		go func() {
			called.Set(int32(waiter.Wait().(ItemState)))
		}()

		time.Sleep(time.Millisecond * 25)
		q.Pop("index")
		time.Sleep(time.Millisecond * 25)
		require.EqualValues(t, ItemPopped, called.Get())
		require.EqualValues(t, 0, q.Len())
	})

	t.Run("fired when cancelled", func(t *testing.T) {
		q := New(FIFO, 3)
		res, waiter := q.AddStateful("item", "index")
		require.True(t, res)

		called := threadsafe.Int32(0)
		go func() {
			called.Set(int32(waiter.Wait().(ItemState)))
		}()

		time.Sleep(time.Millisecond * 25)
		q.CancelAndClose("index")
		time.Sleep(time.Millisecond * 25)
		require.EqualValues(t, ItemCancelled, called.Get())
		require.EqualValues(t, 0, q.Len())
	})
}

func TestLeaks(t *testing.T) {
	t.Run("pop wait", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < 50; i++ {
			wg.Add(1)
			q := New(FIFO, 3)
			go func(q Queue) {
				q.PopWait(DefaultItemIndex).Wait()
				wg.Done()
			}(q)
			go func(q Queue) {
				time.Sleep(time.Millisecond * 25)
				q.Add("test", "")
			}(q)
		}

		wg.Wait()
		goleak.VerifyNone(t)
	})
}
