package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimoutPolicy(t *testing.T) {
	q := New(FIFO, 10, TimeOutPolicy(time.Millisecond*25))
	q.Add("test")
	q.Add("test")
	q.Add("test")
	q.Add("test")
	q.Add("test")
	require.EqualValues(t, "test", q.Pop().(string))
	time.Sleep(time.Millisecond * 50)
	require.Nil(t, q.Pop())
	require.Nil(t, q.Pop())
	require.Nil(t, q.Pop())
	require.Nil(t, q.Pop())
}
