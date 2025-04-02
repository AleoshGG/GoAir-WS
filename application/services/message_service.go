package services

import (
	"GoAir-WS/domain/entities"
	"GoAir-WS/domain/repositories"
	"encoding/json"
	"strconv"
)
type WebSocketHandler interface {
    BroadcastToRoom(roomID string, data []byte)
}

type MessageService struct {
    consumer        repositories.MessageConsumer
    websocketHandler WebSocketHandler
}

func NewMessageService(consumer repositories.MessageConsumer, wsHandler WebSocketHandler) *MessageService {
    return &MessageService{
        consumer:        consumer,
        websocketHandler: wsHandler,
    }
}

func (s *MessageService) StartMessageProcessing() error {
    return s.consumer.ConsumeMessages(func(placeID string, data []byte) {
        var sensor entities.Sensor
        if err := json.Unmarshal(data, &sensor); err == nil {
            s.websocketHandler.BroadcastToRoom(strconv.Itoa(sensor.Id_place), data)
        }
    })
}