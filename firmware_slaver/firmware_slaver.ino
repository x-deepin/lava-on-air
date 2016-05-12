#include <SoftwareSerial.h>
#include <RF24Network.h>
#include <RF24.h>
#include <SPI.h>
#include "protocol.h"


SoftwareSerial mySerial(4,2);
RF24 radio(7,8);

RF24Network network(radio);         
static uint8_t buffer[BufferSize] = {0};

const uint16_t SlaverNode = 02;

void setup_RF(uint16_t node)
{
  SPI.begin();
  radio.begin();
  radio.setPALevel(rf24_pa_dbm_e(RFPALevel));
  radio.setDataRate(rf24_datarate_e(RFDataRate));
  radio.setCRCLength(rf24_crclength_e(RFCRCLen));

  network.begin(RFChannel, node);
}

void setup(void)
{
  Serial.begin(MasterSerialBauds);
  while(!Serial){;}
  
  Serial.println("RF24Network/examples/helloworld_tx/");
  mySerial.begin(9600);
  mySerial.println("Hello 4,5");
  
  setup_RF(SlaverNode);
}


void loop1()
{
  network.update();                          // Check the network regularly

  int c = 0;
  if (c = Serial.available()) {
    int max = c > sizeof(buffer) ? : c;
    Serial.readBytes(buffer, max);
    
    RF24NetworkHeader header(MasterNode, TypeSerialMessage);
    network.write(header, buffer, max);
  }
}

void loop2() 
{
  network.update();                          // Check the network regularly

  int c = 0;
  if (c = mySerial.available()) {
    int max = c > sizeof(buffer) ? : c;
    mySerial.readBytes(buffer, max);
    
    RF24NetworkHeader header(MasterNode, TypeSerialMessage);
    network.write(header, buffer, max);
  }
}

void loop() { loop2(); }
