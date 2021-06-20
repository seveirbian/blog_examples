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

func main() {
	clientID := "bian_publisher_" + strconv.Itoa(time.Now().Second())

	// 1. construct client options
	ops := mqtt.NewClientOptions().SetClientID(clientID).
		AddBroker(common.URL).SetCleanSession(true).SetWill(common.Topic,
			"publisher offline", 1, true)

	// 2. create client and connect to broker
	publisher := mqtt.NewClient(ops)
	token := publisher.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("error connect to broker for %v\n", token.Error())
	}

	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// 3. loop: publish input to broker
	for {
		select {
		case <- sigs:
			publisher.Disconnect(0)
			os.Exit(0)
		default:
			var input string
			fmt.Printf("please input: ")
			fmt.Scanf("%s", &input)
			token := publisher.Publish(common.Topic, 0, false, input)
			if token.WaitTimeout(time.Second*5) && token.Error() != nil {
				fmt.Printf("error publish for %v\n", token.Error())
			}
		}
	}
}
