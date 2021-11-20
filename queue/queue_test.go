package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueueDirection(t *testing.T) {
	t.Run("fifo", func(t *testing.T) {
		q := New(FIFO, 10)

		q.Add("first")
		q.Add("second")
		q.Add("third")

		require.EqualValues(t, "first", q.Pop())
		require.EqualValues(t, "second", q.Pop())
		require.EqualValues(t, "third", q.Pop())
	})

	t.Run("lifo", func(t *testing.T) {
		q := New(LIFO, 10)

		q.Add("first")
		q.Add("second")
		q.Add("third")

		require.EqualValues(t, "third", q.Pop())
		require.EqualValues(t, "second", q.Pop())
		require.EqualValues(t, "first", q.Pop())
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

	require.True(t, q.PopWait().Wait().(bool))
	require.EqualValues(t, 1, q.PopWait().Wait().(int))
	require.EqualValues(t, "t", q.PopWait().Wait().(string))
}
