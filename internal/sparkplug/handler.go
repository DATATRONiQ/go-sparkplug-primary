package sparkplug

import (
	"fmt"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
)

func HandleMessage(msg store.Message) error {
	store.AddMessage(msg)
	switch msg.MessageType {
	case store.NodeBirth:
		fmt.Println("NodeBirth")
	case store.NodeDeath:
		fmt.Println("NodeDeath")
	case store.NodeData:
		fmt.Println("NodeData")
	case store.NodeCommand:
		fmt.Println("NodeCommand")
	case store.DeviceBirth:
		fmt.Println("DeviceBirth")
	case store.DeviceDeath:
		fmt.Println("DeviceDeath")
	case store.DeviceData:
		fmt.Println("DeviceData")
	case store.DeviceCommand:
		fmt.Println("DeviceCommand")
	default:
		return fmt.Errorf("unknown message type: %s", msg.MessageType)
	}
	return nil
}
