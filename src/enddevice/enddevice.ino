/*
 * End-device code for Heltec WiFi LoRa ESP32 v3 board
 * Functions:
 *  - read data from sensors
 *  - parse data to JSON format
 *  - send data through LoRa
 * 
 * Gabriele Crestanello
 * Alma Mater Studiorum, Università di Bologna
*/
#include "LoRaWan_APP.h"
#include "Arduino.h"
#include <ArduinoJson.h>
#include <WiFi.h>
#include "read_data.h"

#define DEEP_SLEEP_INTERVAL                         300 //delay in seconds between data reads


#define RF_FREQUENCY                                868000000 // Hz

#define TX_OUTPUT_POWER                             14        // dBm

#define LORA_BANDWIDTH                              0         // [0: 125 kHz,
                                                              //  1: 250 kHz,
                                                              //  2: 500 kHz,
                                                              //  3: Reserved]
#define LORA_SPREADING_FACTOR                       10         // [SF7..SF12]
#define LORA_CODINGRATE                             1         // [1: 4/5,
                                                              //  2: 4/6,
                                                              //  3: 4/7,
                                                              //  4: 4/8]
#define LORA_PREAMBLE_LENGTH                        10         // Same for Tx and Rx
#define LORA_SYMBOL_TIMEOUT                         0         // Symbols
#define LORA_FIX_LENGTH_PAYLOAD_ON                  false
#define LORA_IQ_INVERSION_ON                        false

#define BUFFER_SIZE                                 128       // LoRa payload size, same for Tx and Rx

#define PIN_HUM 0
#define PIN_TEMP 0


char txpacket[BUFFER_SIZE];
double txNumber;
bool lora_idle=true;
static RadioEvents_t RadioEvents;
char mac[18];

void OnTxDone( void );
void OnTxTimeout( void );
void sendLoRaData(float value, const char* type);
void readHum( void );
void readData( void );

void setup() {
  Serial.begin(115200);
  Mcu.begin();

  txNumber=0;

  uint8_t mac_addr[6];
  WiFi.macAddress(mac_addr);
  snprintf(mac, sizeof(mac), "%02X:%02X:%02X:%02X:%02X:%02X", mac_addr[0], mac_addr[1], mac_addr[2], mac_addr[3], mac_addr[4], mac_addr[5]);

  RadioEvents.TxDone = OnTxDone;
  RadioEvents.TxTimeout = OnTxTimeout;
  
  Radio.Init( &RadioEvents );
  Radio.SetChannel( RF_FREQUENCY );
  Radio.SetTxConfig( MODEM_LORA, TX_OUTPUT_POWER, 0, LORA_BANDWIDTH,
                                  LORA_SPREADING_FACTOR, LORA_CODINGRATE,
                                  LORA_PREAMBLE_LENGTH, LORA_FIX_LENGTH_PAYLOAD_ON,
                                  true, 0, 0, LORA_IQ_INVERSION_ON, 3000 ); 
}


void loop()
{
	if(lora_idle == true)
	{
    readData();
    lora_idle = false;
    esp_sleep_enable_timer_wakeup(DEEP_SLEEP_INTERVAL * 1000000);
    esp_deep_sleep_start();
	}
  Radio.IrqProcess( );
}

void readData(){
  readHum();
  delay(2000);
  readTemp();
}

void readHum(){
  int hum = getHumidity();
  if (hum >= 0) sendLoRaData(hum, "hum");
  else Serial.println("ERROR: humidity read returned negative value, please check sensor connection");
}
void readTemp(){
  // MOCKUP
  sendLoRaData(float(random(200))/2, "temp");
}

void sendLoRaData(float value, const char* type) {
  // random ID generation
  int id = random(100000);
  //int id = txNumber;
  txNumber++;
  // JSON object
  StaticJsonDocument<200> jsonDocument;
  jsonDocument["m"] = mac;
  jsonDocument["i"] = id;
  jsonDocument["v"] = value;
  jsonDocument["t"] = type;

  // serialized JSON to string
  String jsonString;
  serializeJson(jsonDocument, txpacket);

  // LoRa transmission
  Radio.Send( (uint8_t *)txpacket, strlen(txpacket) );

  Serial.print("Dati JSON trasmessi: ");
  Serial.println(txpacket);
}





void OnTxDone( void )
{
	Serial.println("TX done......");
	lora_idle = true;
}
void OnTxTimeout( void )
{
    Radio.Sleep( );
    Serial.println("TX Timeout......");
    lora_idle = true;
}
