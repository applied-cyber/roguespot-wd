package sender

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	scanner "simple-wd/ap-scanner"
)

type APSender struct {
	Queue       [][]scanner.APInfo
	EndpointURL string
	Client      *http.Client
}

func NewSender(endpointURL string) *APSender {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	return &APSender{
		Queue:       [][]scanner.APInfo{},
		EndpointURL: endpointURL,
		Client:      client,
	}
}

func (s *APSender) Send(accessPoints []scanner.APInfo) error {
	log.Printf("Sending %d access points\n", len(accessPoints))

	// Indent so all access points are logged nicely
	accessPointsJson, err := json.MarshalIndent(accessPoints, "", "\t")

	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return err
	}

	log.Printf("%+v\n", string(accessPointsJson))

	req, err := http.NewRequest("POST", s.EndpointURL, bytes.NewBuffer(accessPointsJson))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		s.Queue = append(s.Queue, accessPoints)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP status: %v\n", resp.Status)
		s.Queue = append(s.Queue, accessPoints)
	}
	return nil
}
