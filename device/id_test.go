package device

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntToMAC(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		integer    uint64
		expectedID ID
	}{
		{0, "mac:000000000000"},
		{0x112233445566, "mac:112233445566"},
		{0xF1A293C46570, "mac:f1a293c46570"},
		{0x61FC, "mac:0000000061fc"},
	}

	for _, record := range testData {
		t.Logf("%v", record)
		actualID := IntToMAC(record.integer)
		assert.Equal(record.expectedID, actualID)
	}
}

func TestParseID(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		id           string
		expected     ID
		expectsError bool
	}{
		{"MAC:11:22:33:44:55:66", "mac:112233445566", false},
		{"MAC:11aaBB445566", "mac:11aabb445566", false},
		{"mac:11-aa-BB-44-55-66", "mac:11aabb445566", false},
		{"mac:11,aa,BB,44,55,66", "mac:11aabb445566", false},
		{"uuid:anything Goes!", "uuid:anything Goes!", false},
		{"dns:anything Goes!", "dns:anything Goes!", false},
		{"serial:1234", "serial:1234", false},
		{"mac:11-aa-BB-44-55-66/service", "mac:11aabb445566/service/", false},
		{"mac:11-aa-BB-44-55-66/service/", "mac:11aabb445566/service/", false},
		{"mac:11-aa-BB-44-55-66/service/ignoreMe", "mac:11aabb445566/service/", false},
		{"mac:11-aa-BB-44-55-66/service/foo/bar", "mac:11aabb445566/service/", false},
		{"invalid:a-BB-44-55", "", true},
		{"mac:11-aa-BB-44-55", "", true},
		{"MAC:invalid45566", "", true},
	}

	for _, record := range testData {
		t.Logf("%#v", record)
		id, err := ParseID(record.id)
		assert.Equal(record.expected, id)
		assert.Equal(record.expectsError, err != nil)
		assert.Equal([]byte(record.expected), id.Bytes())
	}
}