package queue

import (
	policies2 "go-threading/queue/policies"
	"testing"
	"time"

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

func TestQueueIndex(t *testing.T) {
	t.Run("pop from non existing index", func(t *testing.T) {
		q := New(FIFO, 10)
		require.Nil(t, q.Pop("non_existing_index"))
	})

	t.Run("default index", func(t *testing.T) {
		q := New(FIFO, 10)
		q.Add("item")
		require.EqualValues(t, "item", q.Pop(DefaultItemIndex))
	})

	t.Run("add new index", func(t *testing.T) {
		q := New(FIFO, 10)
		q.Add("item", "index")
		require.Nil(t, q.Pop(DefaultItemIndex))
		require.EqualValues(t, "item", q.Pop("index"))
		require.Nil(t, q.(*queue).queue["index"])
	})

	t.Run("add multiple new index", func(t *testing.T) {
		q := New(FIFO, 10)
		q.Add("item", "index", "index2", "index3")
		require.Nil(t, q.Pop(DefaultItemIndex))
		require.EqualValues(t, "item", q.Pop("index"))
		require.Nil(t, q.(*queue).queue["index"])
		require.Len(t, q.(*queue).queue["index2"], 1)
		require.Len(t, q.(*queue).queue["index3"], 1)

		require.EqualValues(t, "item", q.Pop("index2"))
		require.Nil(t, q.(*queue).queue["index"])
		require.Nil(t, q.(*queue).queue["index2"])
		require.Len(t, q.(*queue).queue["index3"], 1)

		require.EqualValues(t, "item", q.Pop("index3"))
		require.Nil(t, q.(*queue).queue["index"])
		require.Nil(t, q.(*queue).queue["index2"])
		require.Nil(t, q.(*queue).queue["index3"])
	})
}

func TestEmptyQueue(t *testing.T) {
	q := New(FIFO, 10)
	require.Nil(t, q.Pop(DefaultItemIndex))
}

func TestQueueDirection(t *testing.T) {
	t.Run("fifo", func(t *testing.T) {
		q := New(FIFO, 10)

		q.Add("first")
		q.Add("second")
		q.Add("third")

		require.EqualValues(t, "first", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "second", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "third", q.Pop(DefaultItemIndex))
	})

	t.Run("lifo", func(t *testing.T) {
		q := New(LIFO, 10)

		q.Add("first")
		q.Add("second")
		q.Add("third")

		require.EqualValues(t, "third", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "second", q.Pop(DefaultItemIndex))
		require.EqualValues(t, "first", q.Pop(DefaultItemIndex))
	})
}

func TestPopWait(t *testing.T) {
	q := New(FIFO, 10)

	go func() {
		time.Sleep(time.Millisecond * 60)
		q.Add(true)
		time.Sleep(time.Millisecond * 60)
		q.Add(1)
		time.Sleep(time.Millisecond * 60)
		q.Add("t")
	}()

	require.True(t, q.PopWait(DefaultItemIndex).Wait().(bool))
	require.EqualValues(t, 1, q.PopWait(DefaultItemIndex).Wait().(int))
	require.EqualValues(t, "t", q.PopWait(DefaultItemIndex).Wait().(string))
}

func TestAddWhenFull(t *testing.T) {
	t.Run("multiple indexes > capacity", func(t *testing.T) {
		q := New(FIFO, 3)
		require.False(t, q.Add("item", "index", "index2", "index3", "index4"))
	})
	t.Run("multiple indexes", func(t *testing.T) {
		q := New(FIFO, 3)
		require.True(t, q.Add("item", "index", "index2", "index3"))
		require.False(t, q.Add("item2"))
	})
	t.Run("multiple indexes on immediate eviction", func(t *testing.T) {
		q := New(FIFO, 3, evictImmediatelyF())
		require.True(t, q.Add("item", "index", "index2", "index3"))
		require.True(t, q.Add("item2"))
	})
}
