package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"mqtt_example/common"
)

var messageHandler = func(cli mqtt.Client, mesg mqtt.Message) {
	fmt.Printf("subscriber: message \"%v\" received.\n",
		string(mesg.Payload()))
}

func main() {
	clientID := "bian_subscriber_" + strconv.Itoa(time.Now().Second())

	ops := mqtt.NewClientOptions().SetClientID(clientID).
		AddBroker(common.URL).SetCleanSession(true)

	subscriber := mqtt.NewClient(ops)
	token := subscriber.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("error connect to broker for %v\n", token.Error())
	}

	token = subscriber.Subscribe(common.Topic, 0, messageHandler)
	if token.WaitTimeout(time.Second*5) && token.Error() != nil {
		fmt.Printf("error subscribe for %v\n", token.Error())
	}

	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <- sigs:
			subscriber.Disconnect(0)
			os.Exit(0)
		default:
		}
	}
}
