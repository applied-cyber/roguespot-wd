package sender

import (
	"log"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	scanner "simple-wd/ap-scanner"
)

type APSender struct {
	Queue       []scanner.APInfo
	EndpointURL string
	Client      *http.Client
}

func NewSender(endpointURL string) *APSender {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	return &APSender{
		Queue:       []scanner.APInfo{},
		EndpointURL: endpointURL,
		Client:      client,
	}
}

func (s *APSender) Send(res APInfo) error {
	log.Printf("Sending: %+v\n", res)
	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", s.EndpointURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		s.Queue = append(s.Queue, res)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP status: %v\n", resp.Status)
		s.Queue = append(s.Queue, res)
	}
	return nil
}
