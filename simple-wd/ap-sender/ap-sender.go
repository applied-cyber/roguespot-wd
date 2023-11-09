package sender

import (
	"log"

	scanner "simple-wd/ap-scanner"
)

type APSender struct {
	Queue []scanner.APInfo
}

func NewSender() *APSender {
	return &APSender{
		Queue: []scanner.APInfo{},
	}
}

func (s *APSender) send(res scanner.APInfo) {
	s.Queue = append(s.Queue, res)
	log.Println(s.Queue)
}
