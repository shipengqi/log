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

		assert.Equal(t, expected, len(errs))
	})

	t.Run("output error", func(t *testing.T) {
		opts := NewOptions()
		opts.DisableFile = false
		errs := opts.Validate()

		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "no log output, 'Output' must be set", errs[0].Error())
	})

	t.Run("no enabled logger error", func(t *testing.T) {
		opts := NewOptions()
		opts.DisableConsole = true
		errs := opts.Validate()

		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "no enabled logger, one or more of (DisableConsole, DisableFile) must be set to false", errs[0].Error())
	})

	t.Run("unrecognized level error", func(t *testing.T) {
		opts := NewOptions()
		opts.ConsoleLevel = "errorlevel"
		opts.FileLevel = "errorlevel"
		errs := opts.Validate()

		assert.Equal(t, 2, len(errs))
		assert.Equal(t, "unrecognized level: \"errorlevel\"", errs[0].Error())
		assert.Equal(t, "unrecognized level: \"errorlevel\"", errs[1].Error())
	})
}
