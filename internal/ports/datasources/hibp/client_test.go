package hibp_test

import (
	"encoding/json"
	"gobreach/internal/ports/datasources/hibp"
	_ "gobreach/internal/ports/datasources/hibp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			Name: "happy path",
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mhg := hibp.NewMockHttpGetter(ctrl)

			eEmail := "hacked_account@pwned.com"
			eURL := &url.URL{Scheme: "https", Host: hibp.BaseURL, Path: "breachedaccount/" + eEmail + "?truncateResponse=false"}
			eReq := &http.Request{Method: http.MethodGet, URL: eURL}
			eBody := &hibp.HIBPResponseBody{}
			r := httptest.NewRecorder()
			json.NewEncoder(r.Body).Encode(eBody)
			t.Log(r.Body)
			eResp := &http.Response{StatusCode: http.StatusOK, Body: eReq.Body}

			mhg.EXPECT().Do(eReq).Return(eResp, nil)

			c := hibp.NewClient(mhg, hibp.NewhibpConfig("", ""))
			b, err := c.GetByEmail(eEmail)

			assert.Nil(t, err)
			assert.NotNil(t, b)
			t.Log(b)
		})
	}
}
