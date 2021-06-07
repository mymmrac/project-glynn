package uuid

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	tests := []struct {
		name     string
		uuidStr  string
		expected bool
	}{
		{name: "match", uuidStr: "6d9e6c88-3c82-427e-98b7-7a425e60dfbf", expected: true},
		{name: "match", uuidStr: "9eee574e-2f80-4a07-bb61-238bbcabc239", expected: true},
		{name: "match", uuidStr: "49f6715b-6696-428f-b63a-acf9830bfebb", expected: true},

		{name: "not match", uuidStr: "49f6715b-6696-428f-v63a-acf9830bfebb", expected: false},
		{name: "not match", uuidStr: "49z6715b-6696-428f-v63a-acf9830bfebb", expected: false},
		{name: "not match", uuidStr: "49f6715b-6696-328f-b63a-acf9830bfebb", expected: false},
	}

	rg, err := regexp.Compile(Regex)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := rg.MatchString(tt.uuidStr)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestNew(t *testing.T) {
	rg := regexp.MustCompile(Regex)

	for i := 0; i < 10; i++ {
		uuid := New()
		assert.True(t, rg.MatchString(uuid.String()))
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name        string
		args        args
		expected    UUID
		expectedErr bool
	}{
		{
			name:        "parse ok",
			args:        args{s: "49f6715b-6696-428f-b63a-acf9830bfebb"},
			expected:    UUID{0x49, 0xf6, 0x71, 0x5b, 0x66, 0x96, 0x42, 0x8f, 0xb6, 0x3a, 0xac, 0xf9, 0x83, 0x0b, 0xfe, 0xbb},
			expectedErr: false,
		},
		{
			name:        "parse ok",
			args:        args{s: "F4CDE117-A830-4CF5-A9A5-5C7CF2D9038F"},
			expected:    UUID{0xF4, 0xCD, 0xE1, 0x17, 0xA8, 0x30, 0x4C, 0xF5, 0xA9, 0xA5, 0x5C, 0x7C, 0xF2, 0xD9, 0x03, 0x8F},
			expectedErr: false,
		},
		{
			name:        "parse bad",
			args:        args{s: "F4CDE117-A830-CF5-A9A5-5C7CF2D9038F"},
			expected:    UUID{},
			expectedErr: true,
		},
		{
			name:        "parse bad",
			args:        args{s: "CF5-A9A5-5C7CF2D9038F"},
			expected:    UUID{},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestToStrings(t *testing.T) {
	type args struct {
		uuids []UUID
	}
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name:     "empty",
			args:     args{uuids: []UUID{}},
			expected: []string{},
		},
		{
			name: "one",
			args: args{uuids: []UUID{
				{0xB6, 0x24, 0xDA, 0x99, 0xE1, 0x5D, 0x46, 0xD6, 0xAB, 0x66, 0x41, 0x22, 0x98, 0x3D, 0x42, 0xA1},
			}},
			expected: []string{
				"b624da99-e15d-46d6-ab66-4122983d42a1",
			},
		},
		{
			name: "multiple",
			args: args{uuids: []UUID{
				{0xB6, 0x24, 0xDA, 0x99, 0xE1, 0x5D, 0x46, 0xD6, 0xAB, 0x66, 0x41, 0x22, 0x98, 0x3D, 0x42, 0xA1},
				{0xFC, 0xCA, 0x6E, 0x12, 0x70, 0x20, 0x46, 0xEA, 0xB4, 0x91, 0x97, 0xAD, 0x0D, 0xD8, 0xF1, 0x2D},
				{0xA9, 0x13, 0x2F, 0x92, 0x94, 0x8A, 0x42, 0x46, 0xAE, 0xD3, 0x84, 0xEE, 0x0C, 0x15, 0x62, 0x6F},
			}},
			expected: []string{
				"b624da99-e15d-46d6-ab66-4122983d42a1",
				"fcca6e12-7020-46ea-b491-97ad0dd8f12d",
				"a9132f92-948a-4246-aed3-84ee0c15626f",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToStrings(tt.args.uuids)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
