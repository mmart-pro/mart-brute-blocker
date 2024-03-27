package memorystorage

import (
	"context"
	"testing"

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
	"github.com/stretchr/testify/require"
)

func TestUtils(t *testing.T) {
	t.Run("indexOf", func(t *testing.T) {
		src := []model.Subnet{
			*model.NewSubnet("192.168.1.0/24"),
			*model.NewSubnet("192.168.2.0/24"),
			*model.NewSubnet("192.168.3.0/24"),
		}

		require.Equal(t, indexOf(src, *model.NewSubnet("192.168.0.0/24")), -1)
		require.Equal(t, indexOf(src, *model.NewSubnet("192.168.1.0/24")), 0)
		require.Equal(t, indexOf(src, *model.NewSubnet("192.168.2.0/24")), 1)
		require.Equal(t, indexOf(src, *model.NewSubnet("192.168.3.0/24")), 2)
	})

	t.Run("delete", func(t *testing.T) {
		src := []model.Subnet{
			*model.NewSubnet("192.168.1.0/24"),
			*model.NewSubnet("192.168.2.0/24"),
			*model.NewSubnet("192.168.3.0/24"),
		}

		_, err := deleteSubnet(src, *model.NewSubnet("192.168.0.0/24"))
		require.ErrorIs(t, err, errors.ErrSubnetNotFound)

		remain, err := deleteSubnet(src, *model.NewSubnet("192.168.1.0/24"))
		require.NoError(t, err)
		require.Equal(t, len(src)-1, len(remain))
		require.Equal(t, 0, indexOf(src, *model.NewSubnet("192.168.3.0/24")))
		require.Equal(t, 1, indexOf(src, *model.NewSubnet("192.168.2.0/24")))

		remain, err = deleteSubnet(remain, *model.NewSubnet("192.168.3.0/24"))
		require.NoError(t, err)
		require.Equal(t, len(src)-2, len(remain))
		require.Equal(t, 0, indexOf(src, *model.NewSubnet("192.168.2.0/24")))

		remain, err = deleteSubnet(remain, *model.NewSubnet("192.168.2.0/24"))
		require.NoError(t, err)
		require.Equal(t, len(src)-3, len(remain))
	})

	t.Run("contains", func(t *testing.T) {
		src := []model.Subnet{
			*model.NewSubnet("192.168.1.252/30"),
		}

		require.False(t, contains(src, *model.NewIPAddr("192.168.0.1")))
		require.False(t, contains(src, *model.NewIPAddr("192.168.0.1")))
		require.False(t, contains(src, *model.NewIPAddr("192.168.0.252")))

		require.True(t, contains(src, *model.NewIPAddr("192.168.1.252")))
		require.True(t, contains(src, *model.NewIPAddr("192.168.1.253")))
		require.True(t, contains(src, *model.NewIPAddr("192.168.1.254")))
		require.True(t, contains(src, *model.NewIPAddr("192.168.1.255")))

		src = append(src, *model.NewSubnet("10.58.0.0/16"))
		src = append(src, *model.NewSubnet("172.16.1.0/24"))

		require.False(t, contains(src, *model.NewIPAddr("10.57.0.0")))
		require.False(t, contains(src, *model.NewIPAddr("10.59.0.0")))
		require.False(t, contains(src, *model.NewIPAddr("172.16.0.1")))
		require.False(t, contains(src, *model.NewIPAddr("172.16.2.252")))

		require.True(t, contains(src, *model.NewIPAddr("10.58.1.1")))
		require.True(t, contains(src, *model.NewIPAddr("10.58.255.255")))

		require.True(t, contains(src, *model.NewIPAddr("172.16.1.1")))
		require.True(t, contains(src, *model.NewIPAddr("172.16.1.252")))
	})
}

func TestMemoryStorage(t *testing.T) {
	t.Run("insert exists", func(t *testing.T) {
		ctx := context.Background()
		inmem := NewStorage()

		// InsertBlack
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("10.58.0.0/16")))
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("172.16.1.255/30")))

		// InsertWhite
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("172.16.1.0/30")))
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("127.0.0.0/8")))

		// ExistsBlack
		e, err := inmem.ExistsBlack(ctx, *model.NewSubnet("10.58.0.0/16"))
		require.NoError(t, err)
		require.True(t, e)

		e, err = inmem.ExistsBlack(ctx, *model.NewSubnet("172.16.1.255/30"))
		require.NoError(t, err)
		require.True(t, e)

		e, err = inmem.ExistsBlack(ctx, *model.NewSubnet("172.16.1.0/30"))
		require.NoError(t, err)
		require.False(t, e)

		e, err = inmem.ExistsBlack(ctx, *model.NewSubnet("127.0.0.0/8"))
		require.NoError(t, err)
		require.False(t, e)

		// ExistsWhite
		e, err = inmem.ExistsWhite(ctx, *model.NewSubnet("10.58.0.0/16"))
		require.NoError(t, err)
		require.False(t, e)

		e, err = inmem.ExistsWhite(ctx, *model.NewSubnet("172.16.1.255/30"))
		require.NoError(t, err)
		require.False(t, e)

		e, err = inmem.ExistsWhite(ctx, *model.NewSubnet("172.16.1.0/30"))
		require.NoError(t, err)
		require.True(t, e)

		e, err = inmem.ExistsWhite(ctx, *model.NewSubnet("127.0.0.0/8"))
		require.NoError(t, err)
		require.True(t, e)
	})

	t.Run("contains", func(t *testing.T) {
		ctx := context.Background()
		inmem := NewStorage()

		// InsertBlack
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("10.58.0.0/16")))
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("172.16.1.255/30")))

		// ContainsBlack
		e, err := inmem.ContainsBlack(ctx, *model.NewIPAddr("10.58.1.1"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.252"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.253"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.254"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.255"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.251"))
		require.NoError(t, err)
		require.False(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.1"))
		require.NoError(t, err)
		require.False(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.2"))
		require.NoError(t, err)
		require.False(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.3"))
		require.NoError(t, err)
		require.False(t, e)
		e, err = inmem.ContainsBlack(ctx, *model.NewIPAddr("172.16.1.0"))
		require.NoError(t, err)
		require.False(t, e)

		// InsertWhite
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("172.16.1.0/30")))
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("127.0.0.0/8")))

		// ContainsWhite
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("127.58.1.1"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("172.16.1.0"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("172.16.1.1"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("172.16.1.2"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("172.16.1.3"))
		require.NoError(t, err)
		require.True(t, e)
		e, err = inmem.ContainsWhite(ctx, *model.NewIPAddr("172.16.1.252"))
		require.NoError(t, err)
		require.False(t, e)
	})

	t.Run("delete", func(t *testing.T) {
		ctx := context.Background()
		inmem := NewStorage()

		// InsertBlack
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("10.58.0.0/16")))
		require.NoError(t, inmem.InsertBlack(ctx, *model.NewSubnet("172.16.1.255/30")))

		// InsertWhite
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("172.16.1.0/30")))
		require.NoError(t, inmem.InsertWhite(ctx, *model.NewSubnet("127.0.0.0/8")))

		// Delete with err
		require.ErrorIs(t, errors.ErrSubnetNotFound, inmem.DeleteWhite(ctx, *model.NewSubnet("10.58.0.0/16")))
		require.ErrorIs(t, errors.ErrSubnetNotFound, inmem.DeleteWhite(ctx, *model.NewSubnet("172.16.1.255/30")))

		require.ErrorIs(t, errors.ErrSubnetNotFound, inmem.DeleteBlack(ctx, *model.NewSubnet("172.16.1.0/30")))
		require.ErrorIs(t, errors.ErrSubnetNotFound, inmem.DeleteBlack(ctx, *model.NewSubnet("127.0.0.0/8")))

		// Delete ok
		require.NoError(t, inmem.DeleteBlack(ctx, *model.NewSubnet("10.58.0.0/16")))
		require.NoError(t, inmem.DeleteBlack(ctx, *model.NewSubnet("172.16.1.255/30")))

		require.NoError(t, inmem.DeleteWhite(ctx, *model.NewSubnet("172.16.1.0/30")))
		require.NoError(t, inmem.DeleteWhite(ctx, *model.NewSubnet("127.0.0.0/8")))
	})
}
