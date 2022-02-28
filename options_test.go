package log

import (
	"fmt"
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

func TestOptions_Validate(t *testing.T) {
	opts := NewOptions()
	opts.DisableConsole = true
	opts.DisableFile = true
	errs := opts.Validate()

	expected := `[no enabled logger]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))

	opts.DisableFile = false
	opts.Output = ""
	errs = opts.Validate()

	expected = `[no log output]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))

	opts.ConsoleLevel = "failed"
	opts.FileLevel = "failed"
	errs = opts.Validate()

	expected = `[unrecognized level: "failed" unrecognized level: "failed" no log output]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
