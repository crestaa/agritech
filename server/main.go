package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := MQTT.NewClientOptions().AddBroker(MQTT_BROKER)
	opts.SetClientID("mqtt-subscriber")
	opts.SetUsername(MQTT_USER)
	opts.SetPassword(MQTT_PASS)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Nuovo messaggio ricevuto sul topic: %s - Messaggio: %s\n", msg.Topic(), msg.Payload())
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	topic := MQTT_TOPIC
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("In attesa di messaggi...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect(250)
	fmt.Println("Disconnesso.")
}
