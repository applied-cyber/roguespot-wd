package sender

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type APInfo struct {
	Address string `json:"bssid"`
	SSID    string `json:"ssid"`
}

type APSender struct {
	EndpointURL string
	Client      *http.Client
	Queue       []APInfo
}

func NewSender(endpointURL string) *APSender {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	return &APSender{
		EndpointURL: endpointURL,
		Client:      client,
		Queue:       []APInfo{},
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
