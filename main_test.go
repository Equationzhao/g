package main

import (
	"github.com/Equationzhao/g/internal/config"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/zeebo/assert"
	"io/fs"
	"os"
	"testing"
)

func Test_catchPanic(t *testing.T) {
	gomonkey.ApplyFunc(os.Exit, func(int) {})
	tests := []struct {
		name string
		err  any
	}{
		{
			name: "empty",
			err:  nil,
		},
		{
			name: "ErrExist",
			err:  fs.ErrExist,
		},
		{
			name: "ErrNotExist",
			err:  fs.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catchPanic(tt.err)
		})
	}
}

func Test_preprocessArgs(t *testing.T) {
	settingArgs := []string{"--args2", "--args3"}
	patch := gomonkey.NewPatches()
	defer patch.Reset()
	patch.ApplyFunc(config.Load, func() (*config.Config, error) {
		return &config.Config{
			Args: settingArgs,
		}, nil
	})

	os.Args = []string{"g", "-no-config"}
	preprocessArgs() // this will remove -no-config
	assert.Equal(t, 1, len(os.Args))

	os.Args = []string{"g", "--no-config"}
	preprocessArgs() // this will remove -no-config
	assert.Equal(t, 1, len(os.Args))

	os.Args = []string{"g", "--args1"}
	preprocessArgs() // this will add args from config
	assert.Equal(t, 4, len(os.Args))
	assert.DeepEqual(t, os.Args, []string{"g", "--args2", "--args3", "--args1"})

	os.Args = []string{"g", "--args1", "-no-config"}
	preprocessArgs()
	assert.Equal(t, 2, len(os.Args))

	settingArgs[0] = "-no-config"
	os.Args = []string{"g", "--args1"}
	preprocessArgs()
	assert.Equal(t, 2, len(os.Args))
}
