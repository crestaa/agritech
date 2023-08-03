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

#define RF_FREQUENCY                                868000000 // Hz

#define TX_OUTPUT_POWER                             5        // dBm

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

#define BUFFER_SIZE                                 128 // Define the payload size here

char txpacket[BUFFER_SIZE];

double txNumber;

bool lora_idle=true;

static RadioEvents_t RadioEvents;
void OnTxDone( void );
void OnTxTimeout( void );

void setup() {
  Serial.begin(115200);
  Mcu.begin();

  txNumber=0;

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
    delay(1000);

    sendLoRaData(29.5, "temp");
    lora_idle = false;
	}
  Radio.IrqProcess( );
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



void sendLoRaData(float field, const char* read) {
  // MAC addr + random ID generation
  uint8_t mac[6];
  WiFi.macAddress(mac);
  char macStr[18];
  snprintf(macStr, sizeof(macStr), "%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
  int id = random(100000);

  // JSON object
  StaticJsonDocument<200> jsonDocument;
  jsonDocument["mac"] = macStr;
  jsonDocument["id"] = id;
  jsonDocument["field"] = field;
  jsonDocument["read"] = read;

  // serialized JSON to string
  String jsonString;
  serializeJson(jsonDocument, txpacket);

  // LoRa transmission
  Radio.Send( (uint8_t *)txpacket, strlen(txpacket) );

  Serial.print("Dati JSON trasmessi: ");
  Serial.println(txpacket);
}
