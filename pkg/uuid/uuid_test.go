package uuid

import (
	"reflect"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	tests := []struct {
		name    string
		uuidStr string
		isMatch bool
	}{
		{name: "match", uuidStr: "6d9e6c88-3c82-427e-98b7-7a425e60dfbf", isMatch: true},
		{name: "match", uuidStr: "9eee574e-2f80-4a07-bb61-238bbcabc239", isMatch: true},
		{name: "match", uuidStr: "49f6715b-6696-428f-b63a-acf9830bfebb", isMatch: true},
		{name: "match", uuidStr: "49f6715A-6696-428f-B63A-acf9830bfebb", isMatch: true},
		{name: "match", uuidStr: "F4CDE117-A830-4CF5-A9A5-5C7CF2D9038F", isMatch: true},

		{name: "not match", uuidStr: "49f6715b-6696-428f-v63a-acf9830bfebb", isMatch: false},
		{name: "not match", uuidStr: "49z6715b-6696-428f-v63a-acf9830bfebb", isMatch: false},
		{name: "not match", uuidStr: "49f6715b-6696-328f-b63a-acf9830bfebb", isMatch: false},
	}

	rg, err := regexp.Compile(Regex)
	if err != nil {
		t.Errorf("regexp invalid: %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rg.MatchString(tt.uuidStr); got != tt.isMatch {
				t.Errorf("New() = %v, want %v", got, tt.isMatch)
			}
		})
	}
}

func TestNew(t *testing.T) {
	rg := regexp.MustCompile(Regex)

	for i := 0; i < 10; i++ {
		uuid := New()
		if !rg.MatchString(uuid.String()) {
			t.Errorf("invalid uuid: %q", uuid)
		}
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    UUID
		wantErr bool
	}{
		{
			name:    "parse ok",
			args:    args{s: "49f6715b-6696-428f-b63a-acf9830bfebb"},
			want:    UUID{0x49, 0xf6, 0x71, 0x5b, 0x66, 0x96, 0x42, 0x8f, 0xb6, 0x3a, 0xac, 0xf9, 0x83, 0x0b, 0xfe, 0xbb},
			wantErr: false,
		},
		{
			name:    "parse ok",
			args:    args{s: "F4CDE117-A830-4CF5-A9A5-5C7CF2D9038F"},
			want:    UUID{0xF4, 0xCD, 0xE1, 0x17, 0xA8, 0x30, 0x4C, 0xF5, 0xA9, 0xA5, 0x5C, 0x7C, 0xF2, 0xD9, 0x03, 0x8F},
			wantErr: false,
		},
		{
			name:    "parse bad",
			args:    args{s: "F4CDE117-A830-CF5-A9A5-5C7CF2D9038F"},
			want:    UUID{},
			wantErr: true,
		},
		{
			name:    "parse bad",
			args:    args{s: "CF5-A9A5-5C7CF2D9038F"},
			want:    UUID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStrings(t *testing.T) {
	type args struct {
		uuids []UUID
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{uuids: []UUID{}},
			want: []string{},
		},
		{
			name: "one",
			args: args{uuids: []UUID{
				{0xB6, 0x24, 0xDA, 0x99, 0xE1, 0x5D, 0x46, 0xD6, 0xAB, 0x66, 0x41, 0x22, 0x98, 0x3D, 0x42, 0xA1},
			}},
			want: []string{
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
			want: []string{
				"b624da99-e15d-46d6-ab66-4122983d42a1",
				"fcca6e12-7020-46ea-b491-97ad0dd8f12d",
				"a9132f92-948a-4246-aed3-84ee0c15626f",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStrings(tt.args.uuids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
