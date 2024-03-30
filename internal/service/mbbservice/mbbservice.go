package mbbservice

import (
	"context"
	"errors"
	"time"

	"github.com/mmart-pro/mart-brute-blocker/internal/bucket"
	"github.com/mmart-pro/mart-brute-blocker/internal/config"
	errdef "github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
}

type Storage interface {
	InsertWhite(ctx context.Context, subnet model.Subnet) error
	DeleteWhite(ctx context.Context, subnet model.Subnet) error
	InsertBlack(ctx context.Context, subnet model.Subnet) error
	DeleteBlack(ctx context.Context, subnet model.Subnet) error

	// строгая проверка на совпадение адреса/подсети
	ExistsWhite(ctx context.Context, subnet model.Subnet) (bool, error)
	// строгая проверка на совпадение адреса/подсети
	ExistsBlack(ctx context.Context, subnet model.Subnet) (bool, error)

	// проверка на вхождение адреса/подсети в список
	ContainsWhite(ctx context.Context, ipAddr model.IPAddr) (bool, error)
	// проверка на вхождение адреса/подсети в список
	ContainsBlack(ctx context.Context, ipAddr model.IPAddr) (bool, error)
}

type BucketStorage interface {
	GetBucket(bucketID string) (*bucket.Bucket, error)
	CreateBucket(bucketID string, bucket *bucket.Bucket) error
}

type MbbService struct {
	logger             Logger
	storage            Storage
	ipBucketStorage    BucketStorage
	loginBucketStorage BucketStorage
	pwdBucketStorage   BucketStorage
	svgConfig          config.ServiceConfig
}

func NewMBBService(log Logger, storage Storage,
	ipBucketStorage, loginBucketStorage, pwdBucketStorage BucketStorage, svgConfig config.ServiceConfig) *MbbService {
	return &MbbService{
		logger:             log,
		storage:            storage,
		ipBucketStorage:    ipBucketStorage,
		loginBucketStorage: loginBucketStorage,
		pwdBucketStorage:   pwdBucketStorage,
		svgConfig:          svgConfig,
	}
}

func (service *MbbService) Allow(ctx context.Context, req model.Subnet) error {
	// если уже в белом, то ок
	if exists, err := service.storage.ExistsWhite(ctx, req); err != nil {
		return err
	} else if exists {
		service.logger.Debugf("already allowed %s", req)
		return nil
	}

	// если в чёрном - удалить
	if exists, err := service.storage.ExistsBlack(ctx, req); err != nil {
		return err
	} else if exists {
		if err := service.storage.DeleteBlack(ctx, req); err != nil {
			return err
		}
		service.logger.Debugf("deleted from black list %s", req)
	}

	// добавить в белый
	if err := service.storage.InsertWhite(ctx, req); err != nil {
		return err
	}

	service.logger.Infof("allowed subnet %s", req)
	return nil
}

func (service *MbbService) Deny(ctx context.Context, req model.Subnet) error {
	// если уже в чёрном, то ок
	if exists, err := service.storage.ExistsBlack(ctx, req); err != nil {
		return err
	} else if exists {
		service.logger.Debugf("already denied %s", req)
		return nil
	}

	// если в белом - удалить
	if exists, err := service.storage.ExistsWhite(ctx, req); err != nil {
		return err
	} else if exists {
		if err := service.storage.DeleteWhite(ctx, req); err != nil {
			return err
		}
		service.logger.Debugf("deleted from white list %s", req)
	}

	// добавить в чёрный
	if err := service.storage.InsertBlack(ctx, req); err != nil {
		return err
	}

	service.logger.Infof("denied subnet %s", req)
	return nil
}

func (service *MbbService) Remove(ctx context.Context, req model.Subnet) error {
	service.logger.Debugf("deleting from black list %s", req)
	if err := service.storage.DeleteBlack(ctx, req); err != nil && !errors.Is(err, errdef.ErrSubnetNotFound) {
		return err
	} else if err == nil {
		service.logger.Debugf("deleted from black list %s", req)
		return nil
	}
	service.logger.Debugf("deleting from white list %s", req)
	if err := service.storage.DeleteWhite(ctx, req); err != nil {
		if errors.Is(err, errdef.ErrSubnetNotFound) {
			service.logger.Debugf("%s: %s", err, req)
		}
		return err
	}
	service.logger.Debugf("deleted from white list %s", req)
	return nil
}

func (service *MbbService) Exists(ctx context.Context, req model.Subnet) (model.ListType, error) {
	service.logger.Debugf("white list existence check  %s", req)
	if exists, err := service.storage.ExistsWhite(ctx, req); err != nil {
		return model.NotInList, err
	} else if exists {
		service.logger.Debugf("found in white list %s", req)
		return model.WhiteList, nil
	}

	service.logger.Debugf("black list existence check  %s", req)
	if exists, err := service.storage.ExistsBlack(ctx, req); err != nil {
		return model.NotInList, err
	} else if exists {
		service.logger.Debugf("found in black list %s", req)
		return model.BlackList, nil
	}

	return model.NotInList, nil
}

func (service *MbbService) Contains(ctx context.Context, req model.IPAddr) (model.ListType, error) {
	service.logger.Debugf("white list contains check  %s", req)
	if exists, err := service.storage.ContainsWhite(ctx, req); err != nil {
		return model.NotInList, err
	} else if exists {
		service.logger.Debugf("found in white list %s", req)
		return model.WhiteList, nil
	}

	service.logger.Debugf("black list contains check  %s", req)
	if exists, err := service.storage.ContainsBlack(ctx, req); err != nil {
		return model.NotInList, err
	} else if exists {
		service.logger.Debugf("found in black list %s", req)
		return model.BlackList, nil
	}

	return model.NotInList, nil
}

func (service *MbbService) ClearBucket(_ context.Context, ip model.IPAddr, login string) error {
	if b, err := service.ipBucketStorage.GetBucket(ip.String()); !errors.Is(err, errdef.ErrBucketNotFound) {
		if err != nil {
			return err
		}
		b.Reset()
	}

	if b, err := service.loginBucketStorage.GetBucket(login); !errors.Is(err, errdef.ErrBucketNotFound) {
		if err != nil {
			return err
		}
		b.Reset()
	}

	return nil
}

func (service *MbbService) Check(ctx context.Context, ip model.IPAddr, login, password string) (bool, error) {
	if list, err := service.Contains(ctx, ip); err != nil {
		return false, err
	} else if list == model.WhiteList {
		return true, nil
	} else if list == model.BlackList {
		return false, nil
	}

	if res, err := service.checkBucket(service.ipBucketStorage, ip.String(), service.svgConfig.MaxPerMinForIP); err != nil || !res {
		return false, err
	}

	if res, err := service.checkBucket(service.loginBucketStorage, login, service.svgConfig.MaxPerMinForLogin); err != nil || !res {
		return false, err
	}

	if res, err := service.checkBucket(service.pwdBucketStorage, password, service.svgConfig.MaxPerMinForPassword); err != nil || !res {
		return false, err
	}

	return true, nil
}

func (service *MbbService) checkBucket(bucketStorage BucketStorage, key string, limit int) (bool, error) {
	const capacity = 1

	b, err := bucketStorage.GetBucket(key)
	if err != nil && !errors.Is(err, errdef.ErrBucketNotFound) {
		return false, err
	}

	if errors.Is(err, errdef.ErrBucketNotFound) {
		rate := time.Minute / time.Duration(limit)
		b, err = bucket.NewBucket(capacity, rate)
		if err != nil {
			return false, err
		}
		err = bucketStorage.CreateBucket(key, b)
		if err != nil {
			return false, err
		}
	}
	return b.Allow(), nil
}
