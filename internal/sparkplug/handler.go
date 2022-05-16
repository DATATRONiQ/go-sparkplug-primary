package sparkplug

import (
	"fmt"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/sirupsen/logrus"
)

func HandleMessage(msg store.Message) error {
	store.AddMessage(msg)
	switch msg.MessageType {
	case store.NodeBirth:
		logrus.Debug("NodeBirth")
	case store.NodeDeath:
		logrus.Debug("NodeDeath")
	case store.NodeData:
		logrus.Debug("NodeData")
	case store.NodeCommand:
		logrus.Debug("NodeCommand")
	case store.DeviceBirth:
		logrus.Debug("DeviceBirth")
	case store.DeviceDeath:
		logrus.Debug("DeviceDeath")
	case store.DeviceData:
		logrus.Debug("DeviceData")
	case store.DeviceCommand:
		logrus.Debug("DeviceCommand")
	default:
		return fmt.Errorf("unknown message type: %s", msg.MessageType)
	}
	return nil
}
