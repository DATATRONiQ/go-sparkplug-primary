package sparkplug

import (
	"fmt"
	"strings"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMQTTClient(endpoint, clientID, hostID, user, pass string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(endpoint)
	if user != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}
	opts.SetClientID(clientID)

	stateTopic := fmt.Sprintf("STATE/%s", hostID)
	// as specified in the Sparkplug B Specification
	opts.SetWill(stateTopic, "OFFLINE", 1, true)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("Connected to MQTT broker")
		// as specified in the Sparkplug B Specification
		c.Publish(stateTopic, 1, true, "ONLINE")

		token := c.Subscribe("spBv1.0/+/NBIRTH/+", 1, func(c mqtt.Client, m mqtt.Message) {
			fmt.Println("NBIRTH message received")

			topicParts := strings.Split(m.Topic(), "/")

			store.AddMessage(store.Message{
				ReceivedAt:  time.Now(),
				GroupID:     topicParts[1],
				NodeID:      topicParts[3],
				MessageType: "NBIRTH",
			})
		})
		token.Wait()
		if token.Error() != nil {
			fmt.Println(token.Error())
		}
		fmt.Println("Subscribed to NBIRTH messages")

		token = c.Subscribe("spBv1.0/+/NDEATH/+", 1, func(c mqtt.Client, m mqtt.Message) {
			fmt.Println("NDEATH message received")

			topicParts := strings.Split(m.Topic(), "/")

			store.AddMessage(store.Message{
				ReceivedAt:  time.Now(),
				GroupID:     topicParts[1],
				NodeID:      topicParts[3],
				MessageType: "NDEATH",
			})
		})
		token.Wait()
		if token.Error() != nil {
			fmt.Println(token.Error())
		}
		fmt.Println("Subscribed to NDEATH messages")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
