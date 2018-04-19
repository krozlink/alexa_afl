package betfair

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/krozlink/betting"
	"github.com/valyala/fasthttp"
)

type Credentials struct {
	APIKey         string
	Login          string
	Password       string
	Certificate    []byte
	CertificateKey []byte
}

func GetBetfairSession(c *Credentials) (*betting.Betfair, error) {
	bet := betting.NewBet(c.APIKey)
	key, err := getSessionKey(c.Certificate, c.CertificateKey, c.APIKey, c.Login, c.Password)

	if err != nil {
		return nil, err
	}

	bet.SessionKey = key
	return bet, nil
}

func getSessionKey(certData []byte, keyData []byte, APIKey, login, password string) (string, error) {
	session := &betting.Session{}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return "", err
	}

	client := fasthttp.Client{TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}}

	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.SetRequestURI(betting.CertURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Application", APIKey)
	req.Header.SetMethod("POST")

	bufferString := bytes.NewBuffer([]byte{})
	bufferString.WriteString(`username=`)
	bufferString.WriteString(login)
	bufferString.WriteString(`&password=`)
	bufferString.WriteString(password)

	req.SetBody(bufferString.Bytes())

	err = client.Do(req, resp)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(resp.Body(), session)
	if err != nil {
		return "", err
	}

	switch session.LoginStatus {
	case betting.LS_SUCCESS:
		return session.SessionToken, nil
	default:
		err = errors.New(string(session.LoginStatus))
	}

	return "", err
}
