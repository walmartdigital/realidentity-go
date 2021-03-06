package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type charactersInfo struct {
	Results []results `json:"results"`
}

type results struct {
	Biography biography `json:"biography"`
}

type biography struct {
	FullName string `json:"full-name"`
}

type superheroInfoer interface {
	getName(string) (string, error)
}

type superheroAPIInfoer struct{}

func (s superheroAPIInfoer) getName(superhero string) (string, error) {
	url := fmt.Sprintf("http://superheroapi.com/api/10219222424330451/search/%s", superhero)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error querying the API: %s", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	character := charactersInfo{}
	if err := json.Unmarshal(b, &character); err != nil {
		return "", err
	}
	if len(character.Results) == 0 {
		return "", fmt.Errorf("superhero not found")
	}
	return character.Results[0].Biography.FullName, nil
}

func getSuperheroRealName(s superheroInfoer, superhero string) (string, error) {
	realName, err := s.getName(superhero)
	if err != nil {
		return "", fmt.Errorf("Error querying the superhero API: %s", err)
	}
	return realName, nil
}

func identityHandler(w http.ResponseWriter, r *http.Request) {
	s := r.FormValue("superhero")
	if s == "" {
		http.Error(w, "Missing value", http.StatusBadRequest)
		return
	}
	superheroAPI := superheroAPIInfoer{}
	realname, err := getSuperheroRealName(superheroAPI, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m := struct {
		RealName string `json:"realname"`
	}{RealName: realname}
	msg, _ := json.Marshal(m)
	fmt.Fprintf(w, "%s", msg)
}

func routes() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/identity", identityHandler)
	return r
}

func main() {
	if err := http.ListenAndServe(":8080", routes()); err != nil {
		log.Fatal(err)
	}
}
