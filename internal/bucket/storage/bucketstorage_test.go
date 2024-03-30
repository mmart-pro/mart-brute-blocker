package bucketstorage

import (
	"testing"

	"github.com/mmart-pro/mart-brute-blocker/internal/bucket"
	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestBucketMemoryStorage(t *testing.T) {
	storage := NewBucketMemoryStorage()

	bucket := &bucket.Bucket{}
	bucketID := "test"

	// Create bucket
	err := storage.CreateBucket(bucketID, bucket)
	require.NoError(t, err)

	// Get bucket
	_, err = storage.GetBucket(bucketID)
	require.NoError(t, err)

	// Create existing bucket
	err = storage.CreateBucket(bucketID, bucket)
	require.ErrorIs(t, err, errors.ErrBucketAlreadyExists)

	// Delete bucket
	err = storage.DeleteBucket(bucketID)
	require.NoError(t, err)

	// Test GetBucket with deleted bucket
	_, err = storage.GetBucket(bucketID)
	require.ErrorIs(t, err, errors.ErrBucketNotFound)
}
