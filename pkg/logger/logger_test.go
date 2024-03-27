package logger

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogLevels(t *testing.T) {
	t.Run("test wrong", func(t *testing.T) {
		l, err := NewLogger("wrong", "")
		require.Error(t, err)
		require.Nil(t, l)
	})

	t.Run("test ok", func(t *testing.T) {
		test := []string{"debug", "info", "error", "fatal"}
		for _, tc := range test {
			l, err := NewLogger(tc, "")
			require.NoError(t, err)
			require.NotNil(t, l)
			defer l.Close()
		}
	})
}

func TestLogWrite(t *testing.T) {
	const (
		NotExists string = "NOT EXISTS"
		Exists    string = "MESSAGE LOGGED"
	)

	str, err := captureStderr(func() {
		log, err := NewLogger("info", "")
		require.NoError(t, err)
		defer log.Close()
		log.Debugf(NotExists)
		log.Infof(Exists)
	})
	require.NoError(t, err)

	require.False(t, strings.Contains(str, NotExists))
	require.True(t, strings.Contains(str, Exists))
	require.True(t, strings.Contains(str, `"level":"info"`))
}

func captureStderr(f func()) (string, error) {
	oldStdErr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stderr = writer
	f()
	os.Stderr = oldStdErr
	writer.Close()
	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
