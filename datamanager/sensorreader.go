package datamanager

import (
	"github.com/distributed-go/dto"
)

func SaveReading(reading *dto.SensorMessage) error {
	var sensorreading SensorMessage
	sensorreading.SensorMessage = *reading

	return db.Create(&sensorreading).Error
}
