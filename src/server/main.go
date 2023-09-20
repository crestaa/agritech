package main

import (
	"agritech/server/constants"
	"agritech/server/database"
	"agritech/server/model"
	"database/sql"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var db *sql.DB

func main() {
	// database connection
	db = database.ConnectDB()

	// MQTT broker setup
	opts := MQTT.NewClientOptions().AddBroker(constants.MQTT_BROKER)
	opts.SetClientID("mqtt-subscriber")
	opts.SetUsername(constants.MQTT_USER)
	opts.SetPassword(constants.MQTT_PASS)

	// message handling logic
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		handleMessage(msg)
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	topic := constants.MQTT_TOPIC
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("Waiting for messages...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect(250)
	fmt.Println("Disconnected.")
	defer db.Close()

}

func handleMessage(msg MQTT.Message) {
	data, err := parseJSON(msg)
	if err != nil {
		fmt.Println("ERR during JSON parsing:", err)
		return
	}

	fmt.Println("Received: ", data)

	sensor_id, err := database.GetSensorID(db, data.MAC)
	if err != nil {
		fmt.Println("ERR while getting SensorID:", err)
		return
	}
	measurement_id, err := database.GetMeasurementTypeID(db, data.Type)
	if err != nil {
		fmt.Println("ERR while getting MeasurementTypeID:", err)
		return
	}
	is_double, err := database.CheckDoubles(db, data.ID, sensor_id)
	if err != nil {
		fmt.Println("ERR while checking for doubles:", err)
		return
	}
	if !is_double {
		send := model.Misurazioni{ID_sensore: sensor_id, Nonce: data.ID, Valore: data.Value, ID_tipo_misurazione: measurement_id}
		err = database.SaveMisurazione(db, send)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func parseJSON(msg MQTT.Message) (model.Message, error) {
	var messageData model.Message
	err := json.Unmarshal(msg.Payload(), &messageData)

	return messageData, err
}
