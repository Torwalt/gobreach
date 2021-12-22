package http_test

import (
	"fmt"
	"gobreach/internal/domains/breach"
	_http "gobreach/internal/ports/http"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRoutesGet(t *testing.T) {
	t.Run("TestRoutesGet", func(t *testing.T) {
		domain := "unsafe.com"
		rEmail := fmt.Sprintf("h4cked_mail@%v", domain)
		url := fmt.Sprintf("/breach?email=%v", rEmail)
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			t.Errorf("Error creating a new request: %v", err)
		}

		rr := httptest.NewRecorder()
		ctrl := gomock.NewController(t)
		bSrv := NewMockBreachServer(ctrl)
		b, err := breach.New(rEmail, []string{"Passwords"}, time.Now(), "hacked_site.com")
		bS := []breach.Breach{*b}
		if err != nil {
			t.Errorf("Invalid breach init: %v", err)
		}

		bSrv.EXPECT().GetByEmail(rEmail).Return(bS, nil)
		ro := _http.NewRouter(bSrv, log.Default())

		ro.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestIndex(t *testing.T) {
	t.Run("TestIndex", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)

		if err != nil {
			t.Errorf("Error creating a new request: %v", err)
		}

		rr := httptest.NewRecorder()
		ctrl := gomock.NewController(t)
		bSrv := NewMockBreachServer(ctrl)
		ro := _http.NewRouter(bSrv, log.Default())
		ro.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestRoutesGetNoEmail(t *testing.T) {
	t.Run("TestRoutesGetNoEmail", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/breach", nil)

		if err != nil {
			t.Errorf("Error creating a new request: %v", err)
		}

		rr := httptest.NewRecorder()
		ctrl := gomock.NewController(t)
		bSrv := NewMockBreachServer(ctrl)

		ro := _http.NewRouter(bSrv, log.Default())

		ro.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestRoutesGetBreachServerError(t *testing.T) {
	t.Run("TestRoutesGetNoEmail", func(t *testing.T) {
		domain := "unsafe.com"
		rEmail := fmt.Sprintf("h4cked_mail@%v", domain)
		url := fmt.Sprintf("/breach?email=%v", rEmail)
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			t.Errorf("Error creating a new request: %v", err)
		}

		rr := httptest.NewRecorder()
		ctrl := gomock.NewController(t)
		bSrv := NewMockBreachServer(ctrl)

		bSrv.EXPECT().GetByEmail(rEmail).Return(nil, breach.NewError(breach.DatasourceErr, "TEST datasource err"))
		ro := _http.NewRouter(bSrv, log.Default())

		ro.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

