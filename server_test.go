package main

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	currentDay   = fmt.Sprintf("%v", time.Now().Format("2006-01-02"))
	dateValid    = "2006-01-02"
	dateNotValid = "02.01.2006"
	putJSON      = `{"dateOfBirth":"2006-01-02"}`
	userValid    = "Artem"
	userNotValid = "@rt3m"
)

func TestTodayDate(t *testing.T) {
	valid, today, diff := dateParseCheck(currentDay)
	assert.False(t, valid)
	assert.True(t, today)
	assert.Zero(t, diff)
}

func TestValidDateNotToday(t *testing.T) {
	valid, today, diff := dateParseCheck(dateValid)
	assert.True(t, valid)
	assert.False(t, today)
	assert.NotZero(t, diff)
}

func TestInvalidDateFormat(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	valid, today, diff := dateParseCheck(dateNotValid)
	assert.False(t, valid)
	assert.False(t, today)
	assert.Zero(t, diff)
}

func TestNonLetters(t *testing.T) {
	check := onlyLetters(userNotValid)
	assert.False(t, check)
}

// func TestPutUser(t *testing.T) {
// 	url := "/hello/Artem"
// 	e := echo.New()
// 	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(putJSON))
// 	req.Header.Set("Content-Type", "application/json")
// 	if err != nil {
// 		t.Errorf("The request could not be created because of: %v", err)
// 	}
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	res := rec.Result()
// 	defer res.Body.Close()

// 	if assert.NoError(t, putData(c)) {
// 		assert.Equal(t, http.StatusNoContent, rec.Code)
// 		assert.Equal(t, "", rec.Body.String())
// 	}
// }

// func TestGetUser(t *testing.T) {
// 	url := "/hello/Artem"
// 	e := echo.New()
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	fmt.Printf("\n%v\n", req)
// 	if err != nil {
// 		t.Errorf("The request could not be created because of: %v", err)
// 	}
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	_, _, daysDiff := dateParseCheck("2006-01-02")

// 	res := rec.Result()
// 	defer res.Body.Close()

// 	if assert.NoError(t, getData(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 		assert.Equal(t, fmt.Sprintf("Hello Artem! Your birthday is in %v day(s)", daysDiff), rec.Body.String())
// 	}

// }
