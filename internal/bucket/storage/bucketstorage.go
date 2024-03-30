package bucketstorage

import (
	"sync"

	"github.com/mmart-pro/mart-brute-blocker/internal/bucket"
	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
)

type BucketMemoryStorage struct {
	TokenBuckets sync.Map
}

func NewBucketMemoryStorage() *BucketMemoryStorage {
	bs := &BucketMemoryStorage{}
	return bs
}

func (bs *BucketMemoryStorage) GetBucket(bucketID string) (*bucket.Bucket, error) {
	value, ok := bs.TokenBuckets.Load(bucketID)
	if !ok {
		return nil, errors.ErrBucketNotFound
	}

	b, ok := value.(*bucket.Bucket)
	if !ok {
		return nil, errors.ErrBucketConvertError
	}

	return b, nil
}

func (bs *BucketMemoryStorage) CreateBucket(bucketID string, bucket *bucket.Bucket) error {
	_, err := bs.GetBucket(bucketID)
	if err == nil {
		return errors.ErrBucketAlreadyExists
	}
	bs.TokenBuckets.Store(bucketID, bucket)

	go func(id string, bs *BucketMemoryStorage, shutdown chan bool) {
		<-shutdown
		_ = bs.DeleteBucket(id)
	}(bucketID, bs, bucket.GetDoneChannel())

	return nil
}

func (bs *BucketMemoryStorage) DeleteBucket(bucketID string) error {
	_, err := bs.GetBucket(bucketID)
	if err != nil {
		return err
	}
	bs.TokenBuckets.Delete(bucketID)
	return nil
}
