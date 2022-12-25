package storage

import (
	"context"
	os "os"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
	avalon "github.com/yiwensong/ploggo/avalon"
)

func Test_AvalonJsonStorage_CreatePlayer(t *testing.T) {
	ctx := context.Background()
	tempdir, err := os.MkdirTemp(os.TempDir(), "avalon_storage_test")
	assert.NoError(t, err)

	j, err := NewAvalonJsonStorage(tempdir)
	assert.NoError(t, err)

	player := avalon.NewPlayer("yiwen")

	err = j.CreatePlayer(ctx, player)
	assert.NoError(t, err)
}

func Test_AvalonJsonStorage_GetPlayersById(t *testing.T) {
	ctx := context.Background()
	tempdir, err := os.MkdirTemp(os.TempDir(), "avalon_storage_test")
	assert.NoError(t, err)

	j, err := NewAvalonJsonStorage(tempdir)
	assert.NoError(t, err)

	player := avalon.NewPlayer("yiwen")

	err = j.CreatePlayer(ctx, player)
	assert.NoError(t, err)

	players, err := j.GetPlayersById(ctx, []avalon.PlayerId{player.Id})
	assert.NoError(t, err)

	assert.Equal(t, len(players), 1)
}

func Test_AvalonJsonStorage_StoresStateBetweenLoads(t *testing.T) {
	ctx := context.Background()
	tempdir, err := os.MkdirTemp(os.TempDir(), "avalon_storage_test")
	assert.NoError(t, err)

	j, err := NewAvalonJsonStorage(tempdir)
	assert.NoError(t, err)

	player := avalon.NewPlayer("yiwen")

	err = j.CreatePlayer(ctx, player)
	assert.NoError(t, err)

	err = j.SaveGame(ctx, avalon.NewGame(
		map[avalon.PlayerId]*avalon.PlayerImpl{},
		map[avalon.PlayerId]avalon.Role{},
	))

	j2, err := LoadAvalonJsonStorageFromPath(tempdir)
	assert.NoError(t, err)

	players, err := j2.GetPlayersById(ctx, []avalon.PlayerId{player.Id})
	assert.NoError(t, err)

	assert.Equal(t, len(players), 1)

	assert.Equal(t, len(j2.Games), 1)
}
