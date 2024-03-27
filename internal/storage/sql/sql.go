package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // pgx

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

const (
	whiteListTableName string = "white_list"
	blackListTableName string = "black_list"
)

type SQLStorage struct {
	db  *sql.DB
	dsn string
}

func NewStorage(host string, port int16, user, pwd, db string) *SQLStorage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, db)
	return &SQLStorage{
		dsn: dsn,
	}
}

func (s *SQLStorage) Connect(ctx context.Context) (err error) {
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("connection to db failed: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) InsertWhite(ctx context.Context, subnet model.Subnet) error {
	return s.insertTo(ctx, whiteListTableName, subnet)
}

func (s *SQLStorage) InsertBlack(ctx context.Context, subnet model.Subnet) error {
	return s.insertTo(ctx, blackListTableName, subnet)
}

func (s *SQLStorage) DeleteWhite(ctx context.Context, subnet model.Subnet) error {
	return s.deleteFrom(ctx, whiteListTableName, subnet)
}

func (s *SQLStorage) DeleteBlack(ctx context.Context, subnet model.Subnet) error {
	return s.deleteFrom(ctx, blackListTableName, subnet)
}

func (s *SQLStorage) ExistsWhite(ctx context.Context, subnet model.Subnet) (bool, error) {
	return s.existsIn(ctx, whiteListTableName, subnet)
}

func (s *SQLStorage) ExistsBlack(ctx context.Context, subnet model.Subnet) (bool, error) {
	return s.existsIn(ctx, blackListTableName, subnet)
}

func (s *SQLStorage) ContainsWhite(ctx context.Context, ipAddr model.IPAddr) (bool, error) {
	return s.containsIn(ctx, whiteListTableName, ipAddr)
}

func (s *SQLStorage) ContainsBlack(ctx context.Context, ipAddr model.IPAddr) (bool, error) {
	return s.containsIn(ctx, blackListTableName, ipAddr)
}

func (s *SQLStorage) insertTo(ctx context.Context, dest string, subnet model.Subnet) error {
	q := fmt.Sprintf("insert into %s(subnet) values($1)", dest) //nolint:gosec
	_, err := s.db.ExecContext(ctx, q, subnet.String())
	if err != nil {
		return fmt.Errorf("insert to %s error: %w", dest, err)
	}

	return nil
}

func (s *SQLStorage) deleteFrom(ctx context.Context, dest string, subnet model.Subnet) error {
	q := fmt.Sprintf("delete from %s where subnet = $1", dest) //nolint:gosec
	res, err := s.db.ExecContext(ctx, q, subnet.String())
	if err != nil {
		return fmt.Errorf("delete from %s error: %w", dest, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error: %w", err)
	}
	if rows == 0 {
		return errors.ErrSubnetNotFound
	}

	return nil
}

// строгая проверка на совпадение адреса/подсети
func (s *SQLStorage) existsIn(ctx context.Context, dest string, subnet model.Subnet) (bool, error) {
	q := fmt.Sprintf("select exists(select subnet from %s where subnet = $1 limit 1)", dest) //nolint:gosec
	var exists bool
	err := s.db.QueryRowContext(ctx, q, subnet.String()).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("exists in %s error: %w", dest, err)
	}
	return exists, nil
}

// проверка на вхождение адреса/подсети в список
func (s *SQLStorage) containsIn(ctx context.Context, dest string, ipAddr model.IPAddr) (bool, error) {
	q := fmt.Sprintf("select exists(select subnet from %s where subnet >>= $1 limit 1)", dest) //nolint:gosec
	var exists bool
	err := s.db.QueryRowContext(ctx, q, ipAddr.String()).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("contains in %s error: %w", dest, err)
	}
	return exists, nil
}
