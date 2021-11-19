package queue

import (
	"testing"

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
