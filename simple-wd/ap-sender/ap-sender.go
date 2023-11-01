package sender

import "log"

type APInfo struct {
	Address string
	SSID    string
}

type APSender struct {
	Queue []APInfo
}

func NewSender() *APSender {
	return &APSender{
		Queue: []APInfo{},
	}
}

func (s *APSender) send(res APInfo) {
	s.Queue = append(s.Queue, res)
	log.Println(s.Queue)
}
