package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

type Secret struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUrl                 string `json:"auth_url"`
	TokenUrl                string `json:"token_url"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert"`
	ClientX509CertUrl       string `json:"client_x509_cert"`
}

var (
	App  *firebase.App
	Auth *auth.Client
)

func init() {
	// TODO 環境変数による読み込み

	secret := Secret{
		Type:                    "",
		ProjectId:               "",
		PrivateKeyId:            "",
		PrivateKey:              "",
		ClientEmail:             "",
		ClientId:                "",
		AuthUrl:                 "",
		TokenUrl:                "",
		AuthProviderX509CertUrl: "",
		ClientX509CertUrl:       "",
	}

	file, _ := json.MarshalIndent(secret, "", "")

	ioutil.WriteFile("secret.json", file, 0600)

	opt := option.WithCredentialsFile("./secret.json")

	log.Printf("%+v", opt)

	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	App = app

	Auth, err = app.Auth(context.Background())

	if err != nil {
		log.Fatalf("error create auth client: %v", err)
	}

	os.Remove("./secret.json")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/users", users)

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello World.")
	fmt.Fprintf(w, "Hello World from Go Server.")
}

func login(w http.ResponseWriter, r *http.Request) {

}

func signup(w http.ResponseWriter, r *http.Request) {
	email := ""
	password := ""
	name := ""
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(name)

	u, err := Auth.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("create user error: %v", err)
	}
	fmt.Fprintf(w, u.UID)
}

func users(w http.ResponseWriter, r *http.Request) {
	uid := ""
	u, err := Auth.GetUser(r.Context(), uid)
	if err != nil {
		log.Printf("create user error: %v", err)
	}
	fmt.Fprintf(w, fmt.Sprintf("%+v", u))
}

func authorize(w http.ResponseWriter, r *http.Request) {
	idToken := ""
	token, err := Auth.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
	} else {
		log.Printf("Verified ID token: %v\n", token)
	}
}
