package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/krozlink/alexa_afl/betfair"
	"github.com/krozlink/betting"
	"log"
)

func main() {
	credentials, err := readCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	betfair, err := betfair.GetBetfairSession(credentials)

	getEvents(betfair, "10980856")
	getMarketCatalogues(betfair, "28683484")
	getMarketBooks(betfair, "1.142715718")
}

func readCredentials() (*betfair.Credentials, error) {
	awsConfig := &aws.Config{
		Region: aws.String("ap-southeast-2"),
	}

	sess := session.New(awsConfig)
	apiCh := getParameter(sess, "betfair_api_key")
	loginCh := getParameter(sess, "betfair_login")
	passCh := getParameter(sess, "betfair_password")
	certCh := getParameter(sess, "betfair_certificate")
	keyCh := getParameter(sess, "betfair_private_key")

	api := <-apiCh
	if api.Error != nil {
		return nil, api.Error
	}

	login := <-loginCh
	if login.Error != nil {
		return nil, login.Error
	}

	password := <-passCh
	if password.Error != nil {
		return nil, password.Error
	}

	certificate := <-certCh
	if certificate.Error != nil {
		return nil, certificate.Error
	}

	key := <-keyCh
	if key.Error != nil {
		return nil, key.Error
	}

	credentials := &betfair.Credentials{
		APIKey:         api.Value,
		Certificate:    []byte(certificate.Value),
		CertificateKey: []byte(key.Value),
		Login:          login.Value,
		Password:       password.Value,
	}

	return credentials, nil
}

func getParameter(sess *session.Session, name string) chan ParameterResult {
	sv := ssm.New(sess)
	params := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	}

	result := make(chan ParameterResult)
	go func() {
		resp, err := sv.GetParameter(params)
		if err != nil {
			result <- ParameterResult{
				Value: "",
				Error: err,
			}
		}

		result <- ParameterResult{
			Value: *resp.Parameter.Value,
			Error: nil,
		}
	}()

	return result
}

func getEvents(bet *betting.Betfair, competition string) {

	events, err := bet.ListEvents(betting.Filter{
		Locale: "en",
		MarketFilter: &betting.MarketFilter{
			CompetitionIDs: []string{competition},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range events {
		fmt.Println(printEvent(e.Event))
	}

	fmt.Println(events)
}

func getMarketCatalogues(bet *betting.Betfair, event string) {
	books, err := bet.ListMarketCatalogue(betting.Filter{
		Locale: "en",
		MarketFilter: &betting.MarketFilter{
			EventIDs: []string{event},
		},
		MarketProjection: &[]betting.EMarketProjection{"COMPETITION"},
		MaxResults:       500,
	})

	if err != nil {
		log.Fatalln(err)
	}

	for _, b := range books {
		fmt.Println(printMarketCatalogue(b))
	}

	fmt.Println(books)
}

func getMarketBooks(bet *betting.Betfair, market string) {
	books, err := bet.ListMarketBook(betting.Filter{
		Locale:          "en",
		MarketIDs:       []string{market},
		OrderProjection: "EXECUTABLE",
	})

	if err != nil {
		log.Fatalln(err)
	}

	for _, b := range books {
		fmt.Println(printMarketBook(b))
		fmt.Println(b)
	}
}

func printEvent(e betting.Event) string {
	return fmt.Sprintf("Event Id: %v\nName: %v\nStart Date: %v\n", e.ID, e.Name, e.OpenDate)
}

func printMarketCatalogue(b betting.MarketCatalogue) string {
	return fmt.Sprintf("Market Id: %v\nName: %v\nCompetition: %v\n", b.MarketID, b.MarketName, b.Competition)
}

func printMarketBook(b betting.MarketBook) string {
	return fmt.Sprintf("Market Id: %v\nIs Delayed: %v\nStatus: %v\n", b.MarketID, b.IsMarketDataDelayed, b.Status)
}

type ParameterResult struct {
	Value string
	Error error
}
