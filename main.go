package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/krozlink/betting"
	"log"
	"os"
	"strings"
)

type configuration struct {
	Teams               map[string]string `json:"teams"`
	CompetitionID       string            `json:"competition_id"`
	MatchMarketName     string            `json:"match_market_name"`
	PermiershipMarketID string            `json:"premiership_market_id"`
}

func main() {
	credentials, err := readCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := readConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	betfair, err := NewBetfairSession(credentials)
	if err != nil {
		log.Fatalln(err)
	}

	odds, err := getMatchOdds("Melbourne", config, betfair)
	log.Println(odds)
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

func readConfiguration() (*configuration, error) {
	file, err := os.Open("default_config.json")
	if err != nil {
		return nil, err
	}

	config := &configuration{}

	parser := json.NewDecoder(file)
	if err = parser.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
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

func getMatchOdds(team string, config *configuration, session *betting.Betfair) (float64, error) {
	events := getEvents(session, config.CompetitionID)

	var eventID string
	var isHome bool
	for _, e := range events {
		if strings.HasPrefix(e.Name, team) {
			eventID = e.ID
			isHome = true
		} else if strings.HasSuffix(e.Name, team) {
			eventID = e.ID
			isHome = false
		}
	}

	if eventID == "" {
		return 0, fmt.Errorf("No match found for %v", team)
	}

	marketID, err := getMatchOddsMarket(eventID, config, session)
	if err != nil {
		return 0, err
	}

	runnerIndex := 1
	if isHome {
		runnerIndex = 0
	}
	price, err := getLastPrice(marketID, runnerIndex, session)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func getEvents(bet *betting.Betfair, competition string) []betting.Event {

	result, err := bet.ListEvents(betting.Filter{
		Locale: "en",
		MarketFilter: &betting.MarketFilter{
			CompetitionIDs: []string{competition},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	events := make([]betting.Event, len(result))
	for i, e := range result {
		events[i] = e.Event
	}

	return events
}

func getMatchOddsMarket(event string, config *configuration, session *betting.Betfair) (string, error) {
	books, err := session.ListMarketCatalogue(betting.Filter{
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
		if b.MarketName == config.MatchMarketName {
			return b.MarketID, nil
		}
	}

	return "", fmt.Errorf("No match odds market found")
}

func getLastPrice(market string, index int, session *betting.Betfair) (float64, error) {
	books, err := session.ListMarketBook(betting.Filter{
		Locale:          "en",
		MarketIDs:       []string{market},
		OrderProjection: "EXECUTABLE",
	})

	if err != nil {
		log.Fatalln(err)
	}

	if len(books) != 1 {
		return 0, fmt.Errorf("Unexpected number of books for the market")
	}

	if len(books[0].Runners) != 2 {
		return 0, fmt.Errorf("Unexpected number of runners for the market")
	}

	return books[0].Runners[index].LastPriceTraded, nil
}
