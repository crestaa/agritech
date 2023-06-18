# Agritech

This project is part of my internship at University of Bologna, bachelor degree in Computer Science (a.y. 2022/23).

## Concept
The project is split into:
 - `enddevice`, Arduino sketch(es) used for controlling sensors, collecting data and sending them over LoRaWAN protocol to gateway(s)
 - `gateway`, Arduino sketch(es) for forwarding data from LoRaWAN to the server over MQTT protocol
 - `server`, Go software for collecting data from the MQTT broker, process them and save to the DB
 
 TODO:
 - `DB`
 - `WebApp`

## Usage
### `enddevice` and `gateway`
For these two folders make sure to install all the required libraries. In the `gateway` folder you can edit `credentials.h` file to set up the gateway(s) connection to WiFi and MQTT broker. 

After setting up, you can verify and upload `gateway.ino` and `enddevice.ino` sketches to the boards you're using.

This project is based on Heltec LoRa 32 v3 boards.

### `server`
From inside the `server` folder use `go build` to get yourself an executable file; otherwise use `go run .` to directly execute the code.

Make sure to have Go 1.18 installed.
