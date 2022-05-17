package sparkplug

import (
	"fmt"
	"strings"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/DATATRONiQ/go-sparkplug-primary/third_party/sparkplugb"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func StartMQTTClient(endpoint, clientID, hostID, user, pass string, msgChan chan<- store.Message) {
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
		logrus.Debug("Connected to MQTT broker")
		// as specified in the Sparkplug B Specification
		c.Publish(stateTopic, 1, true, "ONLINE")

		nodeTopics := map[string]byte{
			fmt.Sprintf("spBv1.0/+/%s/+", store.NodeBirth):   1,
			fmt.Sprintf("spBv1.0/+/%s/+", store.NodeDeath):   1,
			fmt.Sprintf("spBv1.0/+/%s/+", store.NodeData):    1,
			fmt.Sprintf("spBv1.0/+/%s/+", store.NodeCommand): 1,
		}

		token := c.SubscribeMultiple(nodeTopics, func(c mqtt.Client, m mqtt.Message) {
			logrus.Debug("node message received")

			if m.Payload() == nil {
				logrus.Warnf("Payload is nil for %s\n", m.Topic())
				return
			}

			topicParts := strings.Split(m.Topic(), "/")

			var payload sparkplugb.Payload
			err := proto.Unmarshal(m.Payload(), &payload)
			if err != nil {
				logrus.Errorf("Failed to unmarshal node message payload of topic %s: %v", m.Topic(), err)
				return
			}

			msgChan <- store.Message{
				ReceivedAt: time.Now(),
				GroupID:    topicParts[1],
				Type:       store.Type(topicParts[2]),
				NodeID:     topicParts[3],
				Payload:    &payload,
			}
		})
		token.Wait()
		if token.Error() != nil {
			logrus.Debug(token.Error())
		}
		logrus.Debug("Subscribed to node messages")

		deviceTopics := map[string]byte{
			fmt.Sprintf("spBv1.0/+/%s/+/+", store.DeviceBirth):   1,
			fmt.Sprintf("spBv1.0/+/%s/+/+", store.DeviceDeath):   1,
			fmt.Sprintf("spBv1.0/+/%s/+/+", store.DeviceData):    1,
			fmt.Sprintf("spBv1.0/+/%s/+/+", store.DeviceCommand): 1,
		}

		token = c.SubscribeMultiple(deviceTopics, func(c mqtt.Client, m mqtt.Message) {
			logrus.Debug("device message received")

			if m.Payload() == nil {
				logrus.Warnf("Payload is nil for %s\n", m.Topic())
				return
			}

			topicParts := strings.Split(m.Topic(), "/")

			var payload sparkplugb.Payload
			err := proto.Unmarshal(m.Payload(), &payload)
			if err != nil {
				logrus.Errorf("Failed to unmarshal node message payload of topic %s: %v", m.Topic(), err)
				return
			}

			msgChan <- store.Message{
				ReceivedAt: time.Now(),
				GroupID:    topicParts[1],
				Type:       store.Type(topicParts[2]),
				NodeID:     topicParts[3],
				DeviceID:   topicParts[4],
				Payload:    &payload,
			}
		})
		token.Wait()
		if token.Error() != nil {
			logrus.Debug(token.Error())
		}
		logrus.Debug("Subscribed to device messages")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
