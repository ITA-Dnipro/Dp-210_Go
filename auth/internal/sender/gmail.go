package sender

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailEmailSender struct {
	CredentialsPath string
	TokenPath       string
	srv             *gmail.Service
}

func NewGmailEmailSender(credPath, tokFile string) (*GmailEmailSender, error) {
	ges := GmailEmailSender{CredentialsPath: credPath, TokenPath: tokFile}
	err := ges.InitializeService()
	if err != nil {
		return nil, err
	}

	return &ges, nil
}

func (ges *GmailEmailSender) InitializeService() error {
	ctx := context.Background()
	b, err := ioutil.ReadFile(ges.CredentialsPath)
	if err != nil {
		return fmt.Errorf("read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		return fmt.Errorf("parse client secret file to config: %w", err)
	}
	client, err := getClient(config, ges.TokenPath)
	if err != nil {
		return fmt.Errorf("get client from config: %w", err)
	}

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("retrieve Gmail client: %w", err)
	}

	ges.srv = srv
	log.Println("gmail email sender successfully initialized")
	return nil
}

func (ges *GmailEmailSender) Send(to, subj, text string) error {
	user := "me"
	var message gmail.Message

	template := fmt.Sprintf("From: 'me'\r\n"+
		"To:  %v\r\n"+
		"Subject: %v \r\n"+
		"\r\n%v", to, subj, text)

	buff := []byte(template)
	message.Raw = base64.StdEncoding.EncodeToString(buff)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err := ges.srv.Users.Messages.Send(user, &message).Do()
	return err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tok, err := tokenFromFile(tokFile)
	if err == nil {
		return config.Client(context.Background(), tok), err
	}
	log.Println(fmt.Errorf("get token from file %v, %w", tokFile, err))

	tok, err = getTokenFromWeb(config)
	if err != nil {
		return nil, fmt.Errorf("get token from web: %w", err)
	}
	if err = saveToken(tokFile, tok); err != nil {
		return nil, fmt.Errorf("save token to file %v: %w", tokFile, err)
	}
	return config.Client(context.Background(), tok), err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("retrieve token from web: %v", err)
	}

	return tok, nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		f, err = os.Create("token.json")
		if err != nil {
			return fmt.Errorf("cache oauth token: %w", err)
		}
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

//Example
// func main() {
// 	ges, err := NewGmailEmailSender("credentials.json", "token.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err = ges.Send("nicknema13@gmail.com", "email test", "gotcha"); err != nil {
// 		log.Fatal(err)
// 	}
// }
