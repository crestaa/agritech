# Agritech

This project is part of my internship at University of Bologna, bachelor degree in Computer Science (a.y. 2022/23).

## Concept
The project is inside the `src` directory and it is split into:
 - `enddevice`, Arduino sketch(es) used for controlling sensors, collecting data and sending them over LoRaWAN protocol to gateway(s)
 - `gateway`, Arduino sketch(es) for forwarding data from LoRaWAN to the server over MQTT protocol
 - `server`, Go software for collecting data from the MQTT broker, process them and save to the DB
 
 TODO:
 - `webapp`

## Setup
### `enddevice` and `gateway`
For these two folders make sure to install all the required libraries. You can edit `./src/gateway/credentials.h` (at first you need to copy it from `./src/gateway/credentials.h.template`) to set up the gateway(s) connection to WiFi and MQTT broker. 

You can edit sensors pins and setup at `./src/enddevice/read_data.h`.

## Usage
### `enddevice` and `gateway`
After setting up, you can compile and upload `gateway.ino` and `enddevice.ino` sketches to the boards you're using.

This project is based on Heltec LoRa 32 v3 boards.

### `server`
From the main folder you can use `docker-compose up -d` to host 3 docker containers:
 - the server, written in Go, available at `./src/server/` and running at port 8080
 - the MQTT broker, which configuration file is available under `./config/mosquitto/` and running at port 1883
 - the MySQL database, running at port 3306

You can edit these configurations inside `docker-compose.yml`.

## Extra
