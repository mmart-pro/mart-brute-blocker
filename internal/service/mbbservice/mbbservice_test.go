package mbbservice

import (
	"context"
	"strconv"
	"testing"
	"time"

	bucketstorage "github.com/mmart-pro/mart-brute-blocker/internal/bucket/storage"
	"github.com/mmart-pro/mart-brute-blocker/internal/config"
	errdef "github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
	memorystorage "github.com/mmart-pro/mart-brute-blocker/internal/storage/memory"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Infof(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func setup() (*MbbService, context.Context) {
	logger := new(MockLogger)
	logger.On("Infof", mock.Anything, mock.Anything).Return()
	logger.On("Debugf", mock.Anything, mock.Anything).Return()

	ipBucketStorage := bucketstorage.NewBucketMemoryStorage()
	loginBucketStorage := bucketstorage.NewBucketMemoryStorage()
	pwdBucketStorage := bucketstorage.NewBucketMemoryStorage()
	cfg := config.ServiceConfig{
		MaxPerMinForLogin:    300,
		MaxPerMinForPassword: 500,
		MaxPerMinForIP:       100,
	}

	service := NewMBBService(logger, memorystorage.NewStorage(), ipBucketStorage, loginBucketStorage, pwdBucketStorage, cfg)
	ctx := context.Background()

	return service, ctx
}

func TestAllowSimple(t *testing.T) {
	service, ctx := setup()

	subnet := model.NewSubnet("192.168.1.0/24")
	err := service.Allow(ctx, *subnet)
	require.NoError(t, err)
}

func TestDenySimple(t *testing.T) {
	service, ctx := setup()

	subnet := model.NewSubnet("192.168.1.0/24")
	err := service.Deny(ctx, *subnet)
	require.NoError(t, err)
}

func TestRemove(t *testing.T) {
	service, ctx := setup()

	subnet := model.NewSubnet("192.168.1.0/24")

	// not in list
	list, err := service.Exists(ctx, *subnet)
	require.NoError(t, err)
	require.Equal(t, model.NotInList, list)

	// remove err
	err = service.Remove(ctx, *subnet)
	require.ErrorIs(t, err, errdef.ErrSubnetNotFound)

	// deny
	err = service.Deny(ctx, *subnet)
	require.NoError(t, err)

	// in black list
	list, err = service.Exists(ctx, *subnet)
	require.NoError(t, err)
	require.Equal(t, model.BlackList, list)

	// remove
	err = service.Remove(ctx, *subnet)
	require.NoError(t, err)

	// not in list
	list, err = service.Exists(ctx, *subnet)
	require.NoError(t, err)
	require.Equal(t, model.NotInList, list)
}

func TestComplex(t *testing.T) {
	service, ctx := setup()

	subnet := model.NewSubnet("192.168.1.0/24")
	err := service.Deny(ctx, *subnet)
	require.NoError(t, err)

	// exists in black list
	list, err := service.Exists(ctx, *subnet)
	require.NoError(t, err)
	require.Equal(t, model.BlackList, list)

	// and ip in black
	ip := model.NewIPAddr("192.168.1.10")
	list, err = service.Contains(ctx, *ip)
	require.NoError(t, err)
	require.Equal(t, model.BlackList, list)

	err = service.Allow(ctx, *subnet)
	require.NoError(t, err)

	// now in white list
	list, err = service.Exists(ctx, *subnet)
	require.NoError(t, err)
	require.Equal(t, model.WhiteList, list)

	// so ip also in white
	ip = model.NewIPAddr("192.168.1.10")
	list, err = service.Contains(ctx, *ip)
	require.NoError(t, err)
	require.Equal(t, model.WhiteList, list)
}

func TestClearBucket(t *testing.T) {
	service, ctx := setup()
	ip := model.NewIPAddr("192.168.1.1")
	const login = "test"

	err := service.ClearBucket(ctx, *ip, login)
	require.NoError(t, err)
}

func TestCheck(t *testing.T) {
	service, ctx := setup()
	ip := model.NewIPAddr("192.168.1.1")
	const login = "test"

	for i := 0; i < 3; i++ {
		allow, err := service.Check(ctx, *ip, login, strconv.Itoa(i))
		require.NoError(t, err)
		require.True(t, allow)
		time.Sleep(time.Minute / time.Duration(service.svgConfig.MaxPerMinForIP))
	}
	allow, err := service.Check(ctx, *ip, login, "password1")
	require.NoError(t, err)
	require.True(t, allow)
	// false, так как 2 запроса подряд
	allow, err = service.Check(ctx, *ip, login, "password2")
	require.NoError(t, err)
	require.False(t, allow)
	// сразу же очистить
	err = service.ClearBucket(ctx, *ip, login)
	require.NoError(t, err)
	// и теперь запрос пройдёт
	allow, err = service.Check(ctx, *ip, login, "password")
	require.NoError(t, err)
	require.True(t, allow)
}
