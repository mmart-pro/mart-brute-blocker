package bucket

import (
	"testing"
	"time"

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/stretchr/testify/require"
)

const rate = time.Second

func TestCreate(t *testing.T) {
	const capacity = 7
	b, err := NewBucket(capacity, rate)
	require.NoError(t, err)

	time.Sleep(rate * 2)
	require.Equal(t, b.Capacity(), b.Amount())
}

func TestBurstCapacity(t *testing.T) {
	const capacity = 42
	b, err := NewBucket(capacity, rate)
	require.NoError(t, err)

	for i := 0; i < capacity-1; i++ {
		require.True(t, b.Allow())
	}
}

func TestEmptyBehavior(t *testing.T) {
	const capacity = 37

	b, err := NewBucket(capacity, rate)
	require.NoError(t, err)

	for b.Amount() > 0 {
		require.True(t, b.Allow())
	}
	// empty - block
	require.Equal(t, false, b.Allow())
}

func TestReset(t *testing.T) {
	const capacity = 37

	b, err := NewBucket(capacity, rate)
	require.NoError(t, err)

	for i := 0; i < capacity/2; i++ {
		require.True(t, b.Allow())
	}
	require.NotEqual(t, b.Amount(), b.Capacity())
	b.Reset()
	require.Equal(t, b.Amount(), b.Capacity())
}

func TestZeroCapacity(t *testing.T) {
	b, err := NewBucket(0, time.Millisecond)
	require.NoError(t, err)

	for i := 0; i < 37; i++ {
		require.False(t, b.Allow())
	}
	require.Equal(t, 0, b.Capacity())
	require.Equal(t, 0, b.Amount())
}

func TestInvalidRate(t *testing.T) {
	_, err := NewBucket(1, 0)
	require.ErrorIs(t, err, errors.ErrBucketRateInvalid)
}

func TestBucketCloseInactive(t *testing.T) {
	b, err := NewBucket(1, rate)
	b.ttl = time.Second * 2
	require.NoError(t, err)
	require.True(t, <-b.GetDoneChannel())
}
