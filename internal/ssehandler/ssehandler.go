package ssehandler

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/api"
	"github.com/sirupsen/logrus"
)

type SSEHandler struct {
	eventChan      chan *api.Event
	newClients     chan chan *api.Event
	closingClients chan chan *api.Event
	clients        map[chan *api.Event]bool
}

func NewSSEHandler() *SSEHandler {
	handler := &SSEHandler{
		eventChan:      make(chan *api.Event, 10),
		newClients:     make(chan chan *api.Event),
		closingClients: make(chan chan *api.Event),
		clients:        make(map[chan *api.Event]bool),
	}
	go handler.start()
	return handler
}

func (sh *SSEHandler) Send(event *api.Event) {
	if event == nil {
		logrus.Error("Received nil event")
		return
	}
	select {
	case sh.eventChan <- event:
	default:
		logrus.Errorf("Dropped event, because msgChan of size %d is full", cap(sh.eventChan))
	}
}

func (sh *SSEHandler) start() {
	for {
		select {
		case nc := <-sh.newClients:
			sh.clients[nc] = true
			logrus.Infof("Added new client to SSEHandler. Clients: %d", len(sh.clients))
		case cc := <-sh.closingClients:
			delete(sh.clients, cc)
			logrus.Infof("Removed client from SSEHandler. Clients: %d", len(sh.clients))
		case event, more := <-sh.eventChan:
			if !more {
				// TODO: Handle clean up
				return
			}
			for clientMsgChan := range sh.clients {
				select {
				case clientMsgChan <- event:
				default:
					logrus.Errorf("Dropped event for client, because clientMsgChan of size %d is full", cap(clientMsgChan))
				}
			}
		}
	}
}

func (sh *SSEHandler) Subscribe(w *bufio.Writer) {
	messageChan := make(chan *api.Event, 10)
	sh.newClients <- messageChan

	defer func() {
		sh.closingClients <- messageChan
	}()

	logrus.Trace("starting stream")

	for event := range messageChan {
		if event == nil {
			logrus.Error("received nil event")
			continue
		}
		bytes, err := json.Marshal(event)
		if err != nil {
			logrus.Errorf("failed to marshal event: %v", err)
			return
		}
		logrus.Debug("sending event")
		_, err = fmt.Fprintf(w, "data: %s\n\n", bytes)
		if err != nil {
			logrus.Errorf("failed to write event: %v", err)
			return
		}
		err = w.Flush()
		if err != nil {
			logrus.Errorf("failed to flush event: %v", err)
			return
		}
		logrus.Debug("flushed event")
	}
	logrus.Debug("stopping stream")
}
