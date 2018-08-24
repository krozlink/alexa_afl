package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/krozlink/betting"
	"log"
)

func main() {
	credentials, err := readCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	betfair, err := GetBetfairSession(credentials)
	if err != nil {
		log.Fatalln(err)
	}

	getEvents(betfair, "10980856")
	getMarketCatalogues(betfair, "28683484")
	getMarketBooks(betfair, "1.142715718")
}

func readCredentials() (*Credentials, error) {
	awsConfig := &aws.Config{
		Region: aws.String("ap-southeast-2"),
	}

	sess := session.New(awsConfig)
	param, err := getParameter(sess, "betfair_credentials")
	if err != nil {
		return nil, err
	}

	credentials := Credentials{}
	if err = json.Unmarshal(param, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

func getParameter(sess *session.Session, name string) ([]byte, error) {
	sv := ssm.New(sess)
	params := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	}
	resp, err := sv.GetParameter(params)
	if err != nil {
		return nil, err
	}

	return []byte(*resp.Parameter.Value), nil
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
