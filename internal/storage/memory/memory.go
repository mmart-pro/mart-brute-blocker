package memorystorage

import (
	"context"
	"sync"

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

type MemoryStorage struct {
	mu sync.Mutex

	blackList []model.Subnet
	whiteList []model.Subnet
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Connect(_ context.Context) error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}

func (s *MemoryStorage) InsertWhite(_ context.Context, subnet model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.whiteList = append(s.whiteList, subnet)
	return nil
}

func (s *MemoryStorage) InsertBlack(_ context.Context, subnet model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.blackList = append(s.blackList, subnet)
	return nil
}

func (s *MemoryStorage) DeleteWhite(_ context.Context, subnet model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list, err := deleteSubnet(s.whiteList, subnet)
	if err == nil {
		s.whiteList = list
	}
	return err
}

func (s *MemoryStorage) DeleteBlack(_ context.Context, subnet model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list, err := deleteSubnet(s.blackList, subnet)
	if err == nil {
		s.blackList = list
	}
	return err
}

// строгая проверка на совпадение адреса/подсети
func (s *MemoryStorage) ExistsWhite(_ context.Context, subnet model.Subnet) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return indexOf(s.whiteList, subnet) >= 0, nil
}

// строгая проверка на совпадение адреса/подсети
func (s *MemoryStorage) ExistsBlack(_ context.Context, subnet model.Subnet) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return indexOf(s.blackList, subnet) >= 0, nil
}

// проверка на вхождение адреса/подсети в список
func (s *MemoryStorage) ContainsWhite(_ context.Context, ipAddr model.IPAddr) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return contains(s.whiteList, ipAddr), nil
}

// проверка на вхождение адреса/подсети в список
func (s *MemoryStorage) ContainsBlack(_ context.Context, ipAddr model.IPAddr) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return contains(s.blackList, ipAddr), nil
}

func deleteSubnet(src []model.Subnet, subnet model.Subnet) ([]model.Subnet, error) {
	i := indexOf(src, subnet)
	if i >= 0 && len(src) > 0 {
		src[i] = src[len(src)-1]
		return src[:len(src)-1], nil
	}
	return nil, errors.ErrSubnetNotFound
}

func indexOf(src []model.Subnet, subnet model.Subnet) int {
	// return slices.IndexFunc(src, func(el model.Subnet) bool { return el.Equal(subnet) })
	for i, v := range src {
		if subnet.Equal(v) {
			return i
		}
	}
	return -1
}

func contains(src []model.Subnet, ipAddr model.IPAddr) bool {
	for _, v := range src {
		if v.Contains(ipAddr) {
			return true
		}
	}
	return false
}
