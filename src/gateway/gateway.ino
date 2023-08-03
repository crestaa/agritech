/*
 * Gateway code for Heltec WiFi LoRa ESP32 v3 board
 * Functions:
 *  - connect to WiFi
 *  - listen over LoRa and receive data
 *  - forward data to MQTT broker
 * 
 * Gabriele Crestanello
 * Alma Mater Studiorum, Università di Bologna
*/
#include "LoRaWan_APP.h"
#include "Arduino.h"
#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include "credentials.h"

#include "Wire.h"
#include "HT_SSD1306Wire.h"

SSD1306Wire  disp(0x3c, 500000, SDA_OLED, SCL_OLED, GEOMETRY_128_64, RST_OLED);


#define RF_FREQUENCY                                868000000 // Hz

#define LORA_BANDWIDTH                              0         // [0: 125 kHz,
                                                              //  1: 250 kHz,
                                                              //  2: 500 kHz,
                                                              //  3: Reserved]
#define LORA_SPREADING_FACTOR                       7         // [SF7..SF12]
#define LORA_CODINGRATE                             1         // [1: 4/5,
                                                              //  2: 4/6,
                                                              //  3: 4/7,
                                                              //  4: 4/8]
#define LORA_PREAMBLE_LENGTH                        8         // Same for Tx and Rx
#define LORA_SYMBOL_TIMEOUT                         0         // Symbols
#define LORA_FIX_LENGTH_PAYLOAD_ON                  false
#define LORA_IQ_INVERSION_ON                        false


#define RX_TIMEOUT_VALUE                            1000
#define BUFFER_SIZE                                 128 // Define the payload size here

char rxpacket[BUFFER_SIZE];

static RadioEvents_t RadioEvents;

int16_t txNumber;

int16_t rssi,rxSize;

bool lora_idle = true;

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

void setup() {
  disp.init();
  disp.setFont(ArialMT_Plain_10);

  Serial.begin(115200);
  Mcu.begin();

  wifi_connect();
  mqtt_connect();
  
  txNumber=0;
  rssi=0;

  RadioEvents.RxDone = OnRxDone;
  Radio.Init( &RadioEvents );
  Radio.SetChannel( RF_FREQUENCY );
  Radio.SetRxConfig( MODEM_LORA, LORA_BANDWIDTH, LORA_SPREADING_FACTOR,
                              LORA_CODINGRATE, 0, LORA_PREAMBLE_LENGTH,
                              LORA_SYMBOL_TIMEOUT, LORA_FIX_LENGTH_PAYLOAD_ON,
                              0, true, 0, 0, LORA_IQ_INVERSION_ON, true );
}

void loop()
{
  if (!mqttClient.connected()) {
    mqtt_connect();
  }
  if(lora_idle)
  {
    lora_idle = false;
    Serial.println("into RX mode");
    Radio.Rx(0);
  }
  Radio.IrqProcess( );
}

void OnRxDone( uint8_t *payload, uint16_t size, int16_t rssi, int8_t snr )
{
  rssi=rssi;
  rxSize=size;
  memcpy(rxpacket, payload, size );
  rxpacket[size]='\0';
  Radio.Sleep();
  Serial.printf("\r\nreceived packet \"%s\" with rssi %d , length %d\r\n",rxpacket,rssi,rxSize);
  mqttClient.publish(MQTT_TOPIC, String(rxpacket).c_str());
  
  disp.clear();
  String s = "received "+String(rxSize)+" bytes";
  disp.drawString(0,0,s);
  s = "with "+String(rssi)+ " rssi";
  disp.drawString(0,10,s);
  disp.display();

  lora_idle = true;
}

void wifi_connect(){
  int attempts = 0;

  Serial.print("Connecting to ");
  Serial.println(WIFI_SSID);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  disp.clear();
  String s = "Connecting to "+String(WIFI_SSID);
  disp.drawString(0,0,s);
  disp.display();

  while (WiFi.status() != WL_CONNECTED && attempts < MAX_ATTEMPTS) {
    delay(1000);
    Serial.print(".");
    attempts++;
  }
  if(attempts >= MAX_ATTEMPTS){
    disp.drawString(0,20,"Can't connect to Wi-Fi.");
    disp.drawString(0,30,"Please check credentials");
    disp.drawString(0,40,"and signal reach.");
    disp.display();
    Serial.println("Can't connect to Wi-Fi. Please check credentials and signal reach.");
    while(1){}
  }

  disp.drawString(0,20,"Connected.");
  disp.display();

  Serial.println("");
  Serial.println("Connected to Wi-Fi");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
  delay(1000);
}

void mqtt_connect(){
  mqttClient.setServer(MQTT_SERVER, MQTT_PORT);
  disp.clear();
  disp.drawString(0,0, "Connecting to MQTT broker");
  disp.display();

  while (!mqttClient.connected()) {
    Serial.println("Trying to connect to MQTT broker...");
    if (mqttClient.connect("HeltecGateway", MQTT_USER, MQTT_PASSWORD)) {
      Serial.println("MQTT broker connected");
      disp.drawString(0,20,"Connected.");
      disp.display();
      delay(1000);
    } else {
      Serial.print("MQTT broker connection failed, error: ");
      Serial.println(mqttClient.state());
      delay(2000);
    }
  }
}