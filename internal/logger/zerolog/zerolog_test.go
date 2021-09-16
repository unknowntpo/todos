package zerolog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintInfo(t *testing.T) {
	out := bytes.NewBufferString("")
	log := New(out)
	log.PrintInfo("test PrintInfo", map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	got := out.String()
	assert.Contains(t, got, "time", "log output must contain field 'time'")
	// See https://github.com/sirupsen/logrus/blob/master/logrus_test.go
	// to understand how to make sure our logger has some fields.
	// Output:
	// {"level":"info","key1":"value1","key2":"value2","message":"test PrintInfo"}
}
