package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// MQTT broker setup
	opts := MQTT.NewClientOptions().AddBroker(MQTT_BROKER)
	opts.SetClientID("mqtt-subscriber")
	opts.SetUsername(MQTT_USER)
	opts.SetPassword(MQTT_PASS)

	// message handling logic
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		handleMessage(msg)
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	topic := MQTT_TOPIC
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	_ = connectDB()

	fmt.Println("Waiting for messages...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect(250)
	fmt.Println("Disconnected.")
}

type Message struct {
	MAC   string  `json:"mac"`
	ID    int     `json:"id"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

func handleMessage(msg MQTT.Message) {
	data, err := parseJSON(msg)
	if err != nil {
		fmt.Println("ERR during JSON parsing:", err)
		return
	}

	fmt.Println("Received: ", data)

}

func parseJSON(msg MQTT.Message) (Message, error) {
	var messageData Message
	err := json.Unmarshal(msg.Payload(), &messageData)

	return messageData, err
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	for err != nil {
		fmt.Println(err)
		time.Sleep(4 * time.Second)
		err = db.Ping()
	}

	fmt.Println("Connessione al database MySQL riuscita!")

	return db
}
