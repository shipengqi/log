package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Default_Options_Validate(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		opts := NewOptions()
		errs := opts.Validate()
		expected := 0

		assert.Equal(t, expected, errs.Len())
	})

	t.Run("output error", func(t *testing.T) {
		opts := NewOptions()
		opts.DisableFile = false
		errs := opts.Validate()

		assert.Equal(t, 1, errs.Len())
		assert.Equal(t, "no log output", errs.Error())
	})

	t.Run("no enabled logger error", func(t *testing.T) {
		opts := NewOptions()
		opts.DisableConsole = true
		errs := opts.Validate()

		assert.Equal(t, 1, errs.Len())
		assert.Equal(t, "no enabled logger", errs.Error())
	})

	t.Run("unrecognized level error", func(t *testing.T) {
		opts := NewOptions()
		opts.ConsoleLevel = "errorlevel"
		opts.FileLevel = "errorlevel"
		errs := opts.Validate()

		assert.Equal(t, 2, errs.Len())
		assert.Equal(t, "unrecognized level: \"errorlevel\" : unrecognized level: \"errorlevel\"", errs.Error())
	})
}
