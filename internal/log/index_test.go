package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := os.CreateTemp("", "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024

	idx, err := newIndex(f, c)
	require.NoError(t, err)

	t.Run("read empty index returns an error", func(t *testing.T) {
		_, _, err := idx.Read(-1)
		require.Error(t, err)
		require.Equal(t, f.Name(), idx.Name())
	})

	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
	}

	t.Run("it writes and reads index without error", func(t *testing.T) {
		for _, want := range entries {
			err := idx.Write(want.Off, want.Pos)
			require.NoError(t, err)

			_, pos, err := idx.Read(int64(want.Off))
			require.NoError(t, err)
			require.Equal(t, want.Pos, pos)
		}
	})

	t.Run("it errors when reading beyond current entries", func(t *testing.T) {
		_, _, err := idx.Read(int64(len(entries)))
		require.Equal(t, io.EOF, err)
		_ = idx.Close()
	})

	t.Run("it builds index state from existing file", func(t *testing.T) {
		f, _ := os.OpenFile(f.Name(), os.O_RDWR, 0600)
		idx, err := newIndex(f, c)
		require.NoError(t, err)

		off, pos, err := idx.Read(-1)
		require.NoError(t, err)
		require.Equal(t, uint32(1), off)
		require.Equal(t, entries[1].Pos, pos)
	})
}