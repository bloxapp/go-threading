package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {

	t.Run("add one", func(t *testing.T) {
		q := New(FIFO, 1000)

		require.Nil(t, q.Pop(DefaultItemIndex))
		require.True(t, q.Add(func() {}))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.Nil(t, q.Pop(DefaultItemIndex))
	})

	t.Run("add multiple", func(t *testing.T) {
		q := New(FIFO, 1000)

		require.True(t, q.Add(func() {}))
		require.True(t, q.Add(func() {}))
		require.True(t, q.Add(func() {}))
		require.True(t, q.Add(func() {}))
		require.True(t, q.Add(func() {}))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.NotNil(t, q.Pop(DefaultItemIndex))
		require.Nil(t, q.Pop(DefaultItemIndex))
	})

}
