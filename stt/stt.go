package stt

import (
	"github.com/mayye4ka/pinder/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Stt struct {
	rabbit *amqp.Connection

	resultChan chan models.SttResult
}

func New(rabbit *amqp.Connection) *Stt {
	return &Stt{
		rabbit:     rabbit,
		resultChan: make(chan models.SttResult, 1024),
	}
}

func (s *Stt) PutTask(task models.SttTask) error {
	s.resultChan <- models.SttResult{
		UserID:    task.UserID,
		MessageID: task.MessageID,
		Text:      "не могу понять. вынь хуй изо рта",
	}
	// TODO:
	return nil
}

func (s *Stt) Start() error {
	// TODO: connect, recv and put to resultsChan
	/*
		t := time.NewTicker(time.Second)
		for range t.C {
			// TODO: somehow load data from rabbit
			var result SttResult

			s.resultChan <- models.SttResult{
				UserID:    result.UserID,
				MessageID: result.MessageID,
				Text:      result.Text,
			}
		}*/
	return nil
}

func (s *Stt) ResultsChan() <-chan models.SttResult {
	return s.resultChan
}
