package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

var redirecturl string

var oauthCfg = &oauth2.Config{
	ClientID:     "22CLZV",
	ClientSecret: "fbb68a421401335786d658f2dd2537ab",
	RedirectURL:  "http://localhost:8080/hello",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.fitbit.com/oauth2/authorize",
		TokenURL: "https://api.fitbit.com/oauth2/token",
	},
	Scopes: []string{"activity"},
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/hello", handleRedirect)
	http.HandleFunc("/authorize", handleAuthorize)
	http.ListenAndServe(":"+os.Getenv("$PORT"), nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	notAuthenticatedTemplate.Execute(w, nil)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	fitbiturl := oauthCfg.AuthCodeURL("")
	fmt.Printf("Url: %v\n", fitbiturl)

	/* u, err := url.Parse(fitbiturl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u) */
	//fmt.Println(u.)
	//fmt.Printf("URL Path: %v\n", r.URL.Path)
	http.Redirect(w, r, fitbiturl, http.StatusFound)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	u, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)

	body := strings.NewReader("clientId=22CLZV&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A5000%2Fhello&code=" + u.Get("code"))
	fmt.Println(body)
	req, err := http.NewRequest("POST", "https://api.fitbit.com/oauth2/token", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Basic MjJDTFpWOmZiYjY4YTQyMTQwMTMzNTc4NmQ2NThmMmRkMjUzN2Fi")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)
	fmt.Println()

	req2, err := http.NewRequest("GET", "https://api.fitbit.com/1/user/-/activities.json", nil)
	if err != nil {
		// handle err
	}
	var BearerString string
	BearerString = "Bearer " + result["access_token"].(string)
	fmt.Printf("Bearer String: %v\n", BearerString)
	req2.Header.Set("Authorization", BearerString)

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		// handle err
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}

	json.NewDecoder(resp2.Body).Decode(&result2)

	fmt.Println(result2)

	fmt.Fprintf(w, "Hello World!!!")
}

var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
	<html><body>
	You have currently not given permissions to access your data. Please authenticate this app with the fitbit provider.
	<form action="/authorize" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
	</body></html>
	`))

var userInfoTemplate = template.Must(template.New("").Parse(`
	<html><body>
	This app is now authenticated to access your Google user info.  Your details are:<br />
	{{.}}
	</body></html>
	`))
