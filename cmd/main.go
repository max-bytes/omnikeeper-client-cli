package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/max-bytes/omnikeeper-client-cli/pkg/credential"
	flag "github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

func main() {
	var keycloakClientID string
	var username string
	var password string
	var omnikeeperURL string
	var graphqlQuery string
	flag.StringVarP(&keycloakClientID, "clientid", "c", "omnikeeper", "Keycloak Client ID")
	flag.StringVarP(&username, "username", "u", "", "Username")
	flag.StringVarP(&password, "password", "p", "", "Password")
	flag.StringVarP(&omnikeeperURL, "omnikeeper-url", "o", "", "Omnikeeper Base URL")
	flag.StringVarP(&graphqlQuery, "query", "q", "", "GraphQL Query")
	flag.BoolP("stdin", "s", false, "Read query from stdin")
	flag.BoolP("suppress-log-output", "i", false, "Suppress log-output, only output response from omnikeeper")
	flag.Parse()

	if !isFlagPassed("omnikeeper-url") {
		log.Fatal("omnikeeper URL not set")
	}

	if isFlagPassed("suppress-log-output") {
		log.SetOutput(ioutil.Discard)
	}

	cs := credential.NewCredentialStore()

	var cred *credentials.Credentials
	if isFlagPassed("username") {
		// username and password is set, we force a login

		// read password from stdin, if not specified
		if !isFlagPassed("password") {
			fmt.Print("Enter Password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatal("Error reading password from stdin", err)
			}
			password = string(bytePassword)
		}

		// store username/password in credential store
		c := &credentials.Credentials{
			ServerURL: omnikeeperURL,
			Username:  username,
			Secret:    password,
		}
		err := cs.Store(c)
		if err != nil {
			log.Printf("Warning: could not store credentials in credential store: %v", err)
		}
		cred = c
	} else {
		// no username specified, rely on stored credentials
		c, err := cs.Get(omnikeeperURL)
		if err != nil {
			log.Fatal("Error getting credentials from credential store", err)
		}

		log.Printf("Re-using stored credentials")
		cred = c
	}

	token, err := login(omnikeeperURL, keycloakClientID, cred.Username, cred.Secret)
	if err != nil {
		log.Fatal("Error logging in", err)
	}
	log.Printf("Successfully logged in to %s", omnikeeperURL)

	var query string
	if isFlagPassed("query") {
		query = graphqlQuery
	} else if isFlagPassed("stdin") {
		// read from stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error reading query from stdin", err)
		}
		query = string(bytes)
	}

	if query != "" {
		graphqlURL := fmt.Sprintf("%s/graphql", omnikeeperURL)
		resp, err := executeGraphql(graphqlURL, token, query)
		if err != nil {
			log.Fatalf("Error executing graphql: %v", err)
		}

		fmt.Print(resp)
	}
}

func login(omnikeeperURL string, keycloakClientID string, username string, password string) (string, error) {

	oAuthEndpoint, err := fetchOAuthInfo(omnikeeperURL)
	if err != nil {
		return "", fmt.Errorf("Error fetching oauth info: %w", err)
	}

	oauth2cfg := &oauth2.Config{
		ClientID: keycloakClientID,
		Endpoint: *oAuthEndpoint,
	}

	// log.Printf("!%s!", username)
	// log.Printf("!%s!", password)
	// log.Print(oAuthEndpoint)
	// log.Print("!")

	ctx := context.Background()
	token, err := oauth2cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("Error getting token: %w", err)
	}

	return token.AccessToken, nil
}

func executeGraphql(graphqlURL string, accessToken string, query string) (string, error) {
	type Payload struct {
		Query string `json:"query"`
	}

	data := Payload{
		Query: query,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", graphqlURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	if resp.StatusCode == http.StatusOK {
		return bodyString, nil
	} else {
		return "", fmt.Errorf("GraphQL query failed with status code %d, returned %s", resp.StatusCode, bodyString)
	}
}

func fetchOAuthInfo(omnikeeperURL string) (*oauth2.Endpoint, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/.well-known/openid-configuration", omnikeeperURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Resp struct {
		TokenEndpoint         string `json:"token_endpoint"`
		AuthorizationEndpoint string `json:"authorization_endpoint"`
	}

	var r Resp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	ret := oauth2.Endpoint{
		AuthURL:  r.AuthorizationEndpoint,
		TokenURL: r.TokenEndpoint,
	}

	return &ret, nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
