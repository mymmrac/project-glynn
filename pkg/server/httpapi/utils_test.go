package httpapi

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_respondJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	type args struct {
		data       interface{}
		statusCode int
	}
	type expected struct {
		data string
		err  bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "ok",
			args: args{
				data: struct {
					Test string `json:"test"`
					Ok   bool   `json:"ok"`
				}{"test", true},
				statusCode: http.StatusOK,
			},
			expected: expected{
				data: "{\"test\":\"test\",\"ok\":true}\n",
				err:  false,
			},
		},
		{
			name: "nil",
			args: args{
				data:       nil,
				statusCode: http.StatusOK,
			},
			expected: expected{
				data: "",
				err:  true,
			},
		},
		{
			name: "err",
			args: args{
				data:       func() {},
				statusCode: http.StatusOK,
			},
			expected: expected{
				data: "",
				err:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := respondJSON(rr, tt.args.data, tt.args.statusCode)
			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.data, rr.Body.String())
			assert.Equal(t, tt.args.statusCode, rr.Code)
			assert.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
		})
	}
}

func Test_respondJSONError(t *testing.T) {
	rr := httptest.NewRecorder()

	type args struct {
		error      error
		statusCode int
	}
	type expected struct {
		data string
		err  bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "ok",
			args: args{
				error:      errors.New("test"),
				statusCode: http.StatusBadRequest,
			},
			expected: expected{
				data: "{\"error\":\"test\"}\n",
				err:  false,
			},
		},
		{
			name: "nil",
			args: args{
				error:      nil,
				statusCode: http.StatusInternalServerError,
			},
			expected: expected{
				data: "",
				err:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := respondJSONError(rr, tt.args.error, tt.args.statusCode)
			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.data, rr.Body.String())
			assert.Equal(t, tt.args.statusCode, rr.Code)
		})
	}
}

func Test_decodeJSON(t *testing.T) {
	type testData struct {
		Test string
	}
	type args struct {
		data string
	}
	type expected struct {
		v   *testData
		err bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "ok",
			args: args{
				data: "{\"test\":\"test\"}\n",
			},
			expected: expected{
				v:   &testData{Test: "test"},
				err: false,
			},
		},
		{
			name: "err",
			args: args{
				data: "{",
			},
			expected: expected{
				v:   nil,
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.data))

			var actual *testData
			err := decodeJSON(r, &actual)
			if tt.expected.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.v, actual)
		})
	}
}
