package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type fakeInfoer struct {
	realName string
	err      error
}

func (f fakeInfoer) getName(name string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.realName, nil
}

func TestGetSuperheroRealName(t *testing.T) {
	tt := []struct {
		tn            string
		fake          fakeInfoer
		expectedName  string
		expectedError error
	}{
		{
			tn: "superhero doesn't exist",
			fake: fakeInfoer{
				realName: "",
				err:      errors.New("superhero not found"),
			},
			expectedError: errors.New("Error querying the superhero API: superhero not found"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tn, func(t *testing.T) {
			_, err := getSuperheroRealName(tc.fake, "")
			if !reflect.DeepEqual(err, tc.expectedError) {
				t.Errorf("Expected error to be: %v but got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestIdentityHandler(t *testing.T) {
	tt := []struct {
		tn            string
		superhero     string
		expectedName  string
		expectedError error
	}{
		{tn: "real name is wade wilson", superhero: "deadpool", expectedName: `{"realname":"Wade Wilson"}`, expectedError: nil},
		{tn: "missing value", superhero: "", expectedName: "", expectedError: errors.New("")},
		{tn: "superhero not found", superhero: "nohero", expectedName: "", expectedError: errors.New("")},
	}
	for _, tc := range tt {
		t.Run(tc.tn, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8080/identity?superhero="+tc.superhero, nil)
			if err != nil {
				t.Fatalf("Could not create the request: %v", err)
			}
			recorder := httptest.NewRecorder()
			identityHandler(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Could not read the response: %v", err)
			}
			if !reflect.DeepEqual(err, tc.expectedError) {
				if res.StatusCode != http.StatusBadRequest {
					t.Fatalf("Expected status Bad Request, but got %v", res.StatusCode)
				}
				return
			}
			name := string(bytes.TrimSpace(b))
			if name != tc.expectedName {
				t.Fatalf("Expected name to be: %s, but got: %s", tc.expectedName, name)
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(routes())
	defer srv.Close()

	url := fmt.Sprintf("%s%s?superhero=deadpool", srv.URL, "/identity")
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Could not send the request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, but got: %s", resp.Status)
	}
}
