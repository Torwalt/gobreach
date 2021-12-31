package hibp_test

import (
	"bytes"
	"fmt"
	"gobreach/cmd/server/config"
	"gobreach/internal/domains/breach"
	"gobreach/internal/ports/datasources/hibp"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	testRespSuccess = `[{"Name":"Adobe",
"Title":"Adobe",
"Domain":"adobe.com",
"BreachDate":"2013-10-04",
"AddedDate":"2013-12-04T00:00Z",
"ModifiedDate":"2013-12-04T00:00Z",
"PwnCount":152445165,
"Description":"In October 2013, 153 million Adobe accounts were breached with each containing an internal ID, username, email, <em>encrypted</em> password and a password hint in plain text. The password cryptography was poorly done and <a href=\"http://stricture-group.com/files/adobe-top100.txt\" target=\"_blank\" rel=\"noopener\">many were quickly resolved back to plain text</a>. The unencrypted hints also <a href=\"http://www.troyhunt.com/2013/11/adobe-credentials-and-serious.html\" target=\"_blank\" rel=\"noopener\">disclosed much about the passwords</a> adding further to the risk that hundreds of millions of Adobe customers already faced.",
"DataClasses":["Email addresses","Password hints","Passwords","Usernames"],
"IsVerified":true,
"IsSensitive":false,
"IsRetired":false,
"IsSpamList":false}]`
	testMalformedResp = `[{"Name":"Adobe",
"Title":"Adobe",
"Domain":"adobe.com",
"BreachDate":"2013/10/04",
"AddedDate":"2013-12-04T00:00Z",
"ModifiedDate":"2013-12-04T00:00Z",
"PwnCount":152445165,
"Description":"In October 2013, 153 million Adobe accounts were breached with each containing an internal ID, username, email, <em>encrypted</em> password and a password hint in plain text. The password cryptography was poorly done and <a href=\"http://stricture-group.com/files/adobe-top100.txt\" target=\"_blank\" rel=\"noopener\">many were quickly resolved back to plain text</a>. The unencrypted hints also <a href=\"http://www.troyhunt.com/2013/11/adobe-credentials-and-serious.html\" target=\"_blank\" rel=\"noopener\">disclosed much about the passwords</a> adding further to the risk that hundreds of millions of Adobe customers already faced.",
"DataClasses":["Email addresses","Password hints","Passwords","Usernames"],
"IsVerified":true,
"IsSensitive":false,
"IsRetired":false,
"IsSpamList":false}]`
	testEmptyResp = `{}`
	test429Resp   = `{"statusCode": 429, "message": "Rate limit is exceeded. Try again in 6 seconds."}`
)

type MockSleepCounter uint8

func (m *MockSleepCounter) Sleep(d time.Duration) {
	*m++
}

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		Name string
		Data []byte
	}{
		{
			Name: "success",
			Data: []byte(testRespSuccess),
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mhg := hibp.NewMockHttpGetter(ctrl)

			eEmail := "hacked_account@pwned.com"
			r := ioutil.NopCloser(bytes.NewReader(test.Data))
			eResp := &http.Response{StatusCode: http.StatusOK, Body: r}
			mhg.EXPECT().Do(gomock.Any()).Return(eResp, nil)

			var msc MockSleepCounter
			c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
			bS, err := c.GetByEmail(eEmail)

			assert.Equal(t, len(bS), 1)
			assert.Nil(t, err)
			b := bS[0]
			assert.Equal(t, "Adobe", b.BreachSource)
			eHash := breach.HashEmail(eEmail)
			assert.Equal(t, eHash, b.EmailHash)
			assert.Equal(t, "2013-10-04", b.BreachDate.Format("2006-01-02"))
		})
	}
}

func TestGetByEmailRequestError(t *testing.T) {
	t.Parallel()
	t.Run("error resp", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mhg := hibp.NewMockHttpGetter(ctrl)

		eEmail := "hacked_account@pwned.com"
		mhg.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("Test Error"))

		var msc MockSleepCounter
		c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
		b, err := c.GetByEmail(eEmail)

		assert.NotNil(t, err)
		assert.Nil(t, b)
	})
}

func TestGetByEmailDecodeError(t *testing.T) {
	t.Parallel()
	t.Run("decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mhg := hibp.NewMockHttpGetter(ctrl)
		data := []byte(`not json`)

		eEmail := "hacked_account@pwned.com"
		r := ioutil.NopCloser(bytes.NewReader(data))
		eResp := &http.Response{StatusCode: http.StatusOK, Body: r}
		mhg.EXPECT().Do(gomock.Any()).Return(eResp, nil)

		var msc MockSleepCounter
		c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
		b, err := c.GetByEmail(eEmail)

		assert.NotNil(t, err)
		assert.Nil(t, b)
	})
}

func TestGetByEmailToBreachError(t *testing.T) {
	t.Parallel()
	t.Run("breach error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mhg := hibp.NewMockHttpGetter(ctrl)

		eEmail := "hacked_account@pwned.com"
		data := []byte(testMalformedResp)
		r := ioutil.NopCloser(bytes.NewReader(data))
		eResp := &http.Response{StatusCode: http.StatusOK, Body: r}
		mhg.EXPECT().Do(gomock.Any()).Return(eResp, nil)

		var msc MockSleepCounter
		c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
		b, err := c.GetByEmail(eEmail)

		assert.NotNil(t, err)
		assert.Nil(t, b)
	})
}

func TestGetByEmailNot200(t *testing.T) {
	tests := []struct {
		Name       string
		Data       []byte
		StatusCode int
	}{
		{
			Name:       "no breach found",
			Data:       []byte(testEmptyResp),
			StatusCode: http.StatusNotFound,
		},
		{
			Name:       "not 200 resp",
			Data:       []byte(testEmptyResp),
			StatusCode: http.StatusServiceUnavailable,
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mhg := hibp.NewMockHttpGetter(ctrl)

			eEmail := "hacked_account@pwned.com"
			r := ioutil.NopCloser(bytes.NewReader(test.Data))
			eResp := &http.Response{StatusCode: test.StatusCode, Body: r}
			mhg.EXPECT().Do(gomock.Any()).Return(eResp, nil)

			var msc MockSleepCounter
			c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
			b, err := c.GetByEmail(eEmail)

			assert.NotNil(t, err)
			assert.Nil(t, b)
		})
	}
}

func TestGetByEmailRetry(t *testing.T) {
	t.Run("TestGetByEmailRetry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mhg := hibp.NewMockHttpGetter(ctrl)
		eEmail := "hacked_account@pwned.com"

		// 429 resp
		r429 := ioutil.NopCloser(bytes.NewReader([]byte(test429Resp)))
		eResp429 := &http.Response{StatusCode: http.StatusTooManyRequests, Body: r429}
		first := mhg.EXPECT().Do(gomock.Any()).Return(eResp429, nil)

		// 200 resp
		r200 := ioutil.NopCloser(bytes.NewReader([]byte(testRespSuccess)))
		eResp200 := &http.Response{StatusCode: http.StatusOK, Body: r200}
		second := mhg.EXPECT().Do(gomock.Any()).Return(eResp200, nil)

		gomock.InOrder(first, second)

		var msc MockSleepCounter
		c := hibp.NewClient(mhg, hibp.NewhibpConfig("", "", 2), msc.Sleep)
		bS, err := c.GetByEmail(eEmail)

		assert.Equal(t, len(bS), 1)
		assert.Nil(t, err)
		b := bS[0]
		assert.Equal(t, "Adobe", b.BreachSource)
		eHash := breach.HashEmail(eEmail)
		assert.Equal(t, eHash, b.EmailHash)
		assert.Equal(t, "2013-10-04", b.BreachDate.Format("2006-01-02"))
		assert.Equal(t, int(msc), 1)
	})
}

/* integraion tests */

func TestIntegrationGetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := config.FromEnv()
	hconf := hibp.NewhibpConfig(c.HIBPHost, c.HIBPKey, 2)
	// we should use the real sleep, otherwise the test may fail
	hclient := hibp.NewClient(http.DefaultClient, hconf, time.Sleep)
	// https://haveibeenpwned.com/API/v3#TestAccounts
	b, err := hclient.GetByEmail("account-exists@hibp-integration-tests.com")

	if err != nil {
		t.Errorf(err.Error())
	}

	assert.NotNil(t, b)
	t.Log(b)

}
