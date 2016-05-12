#include <SoftwareSerial.h>
#include <RF24Network.h>
#include <RF24.h>
#include <SPI.h>
#include "protocol.h"

const uint16_t SlaverNode = 02;

SoftwareSerial mySerial(4,2);
RF24 radio(PIN_RF_CE, PIN_RF_CSN);

RF24Network network(radio);         
static uint8_t buffer[BufferSize] = {0};

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

static uint8_t fpanel_pin = 0;
static unsigned long fpanel_endtime = 0;

void handle_fpanel_control(const FPanelControl cmd)
{
  if (fpanel_pin != 0) {
    // TODO: Report interrupt 
    digitalWrite(fpanel_pin, LOW);
  }
  
  fpanel_pin = cmd.pin;
  fpanel_endtime = millis() + cmd.duration;
  
  digitalWrite(fpanel_pin, LOW);
  delay(10);
  digitalWrite(fpanel_pin, HIGH);
}
void fpanel_control_loop()
{
  if (fpanel_pin == 0 || fpanel_endtime == 0) {
    return;
  }

  if (fpanel_endtime > millis()) {
    digitalWrite(fpanel_pin, LOW);
    fpanel_pin = 0;
    fpanel_endtime = 0;
  }
}

void fpanel_status_loop()
{
  static unsigned long time = millis();
  if (millis() - time  > 2000) {
    time = millis();
    RF24NetworkHeader header(MasterNode, TypeFPanelStatus);

    uint8_t status = 0;
    if (digitalRead(PIN_POWER_LED)) {
      status +=   FPanelStatusPower;
    }
    network.write(header, &status, 1);
  }
}


void forward_out_serial_message()
{
  uint8_t available = 0;
  uint8_t rd = 0;
  
  while ((available = mySerial.available()) > 0 && rd < BufferSize) {
    
    rd += mySerial.readBytes(buffer + rd,
			     available > (BufferSize - rd) ?
			     (BufferSize - rd) :
			     available
			     );
  }
  
  if (rd != 0) {
    RF24NetworkHeader header(MasterNode, TypeSerialMessage);
    network.write(header, buffer, rd);
    network.update();
  }
}

void loop()
{
  forward_out_serial_message();
  
  network.update();

  if ( network.available() )  {
    RF24NetworkHeader h;
    uint8_t c = network.read(h, buffer, BufferSize);

    switch (h.type) {
    case TypeSerialMessage:
      Serial.write(buffer, c);
      break;
    case TypeFPanelControl:
      if (c == sizeof(FPanelControl)) {
	FPanelControl cmd;
	memcpy(&cmd, buffer, sizeof(FPanelControl));
	handle_fpanel_control(cmd);
      } else  {
	// Handle invalid message
      }
      break;
    }
  }

  fpanel_control_loop();
  fpanel_status_loop();
}
