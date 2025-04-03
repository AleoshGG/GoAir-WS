package services

import (
	"encoding/json"
	"fmt"
	"log"

	"GoAir-WS/domain/entities"
	"GoAir-WS/domain/repositories"
)

type MessageService struct {
	messageRepo repositories.MessageRepository
	wsAdapter   WebSocketAdapter
}

type WebSocketAdapter interface {
	BroadcastSensor(sensor entities.Sensor)
	BroadcastUserRequest(req entities.UserRequest)
	BroadcastConfirmationInstallation(msg entities.ConfirmInstalltionMessage)
}

func NewMessageService(msgRepo repositories.MessageRepository, wsAdapter WebSocketAdapter) *MessageService {
	return &MessageService{
		messageRepo: msgRepo,
		wsAdapter:   wsAdapter,
	}
}

// StartMessageProcessing inicia dos goroutines: una para mensajes Sensor y otra para UserRequest.
func (s *MessageService) StartMessageProcessing() {
	sensorMsgs := s.messageRepo.ConsumeSensorMessages()
	userReqMsgs := s.messageRepo.ConsumeUserRequestMessages()
	cInstallation := s.messageRepo.ConsumeConfirmInstallationMessages()

	go func() {
		for msg := range sensorMsgs {
			var sensor entities.Sensor
			if err := json.Unmarshal(msg.Body, &sensor); err != nil {
				log.Println("Error decoding sensor message:", err)
				continue
			}
			// Se espera que los clientes se registren usando el IdPlace (convertido a string)
			s.wsAdapter.BroadcastSensor(sensor)
		}
	}()

	go func() {
		for msg := range userReqMsgs {
			var req entities.UserRequest
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				log.Println("Error decoding user request message:", err)
				continue
			}
			// Se espera que los clientes se registren usando el Destination
			s.wsAdapter.BroadcastUserRequest(req)
		}
	}()

	go func() {
		for msg := range cInstallation {
			var cIMsg entities.ConfirmInstalltionMessage
			if err := json.Unmarshal(msg.Body, &cIMsg); err != nil {
				log.Println("Error decoding user request message:", err)
				continue
			}
			fmt.Print(cIMsg)
			// Se espera que los clientes se registren usando el Destination
			s.wsAdapter.BroadcastConfirmationInstallation(cIMsg)
		}
	}()
}
