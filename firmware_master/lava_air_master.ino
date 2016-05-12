#include <SPI.h>
#include <RF24.h>
#include <RF24Network.h>

#include "protocol.h"


static uint8_t buffer[BufferSize] = {0};

RF24 radio(PIN_RF_CE, PIN_RF_CSN);                // nRF24L01(+) radio attached using Getting Started board 

RF24Network network(radio);      // Network uses that radio


void setup_RF(uint16_t node)
{
  SPI.begin();
  radio.begin();
  radio.setPALevel(rf24_pa_dbm_e(RFPALevel));
  radio.setDataRate(rf24_datarate_e(RFDataRate));
  radio.setCRCLength(rf24_crclength_e(RFCRCLen));

  network.begin(RFChannel, node);
}

void send_message(uint16_t from, uint16_t to, uint8_t type, uint16_t id, const void* data, uint8_t data_len)
{
  Header h = {
    .from_node = from,
    .to_node = to,
    .id = id,
    .type = type,
    .payload_size = data_len,
  };
  Serial.write((uint8_t*)&h, sizeof(Header));
  
  if (data_len != 0) {
    Serial.write((uint8_t*)data, data_len);
    Serial.flush();
  }
}


inline
void send_hello()
{
  RF24NetworkHeader  h(MasterNode, TypeHello);
  send_message(h.from_node, h.to_node, h.type, h.id, "hello", sizeof("hello"));
}

void setup(void)
{
  Serial.begin(MasterSerialBauds);
  
  setup_RF(MasterNode);

  send_hello();
}
 

bool handle_serial_line()
{
  uint16_t to = Serial.read();
  to = (to << 8) + Serial.read();
  if (!network.is_valid_address(to)) {
        return false;
  }
  
  uint8_t type = Serial.read();
  if (CheckTypeValidation((T)type)) {
        return false;
  }

  uint8_t size = Serial.read();
  for (uint8_t i=0; i < size; i++) {
    buffer[i] = Serial.read();
  }
  RF24NetworkHeader header(to, type);
  network.write(header, buffer, size);
  return true;
}


inline
void discard_serial_line_buffer()
{
  while (Serial.available()) {
    Serial.read();
  }
  
  RF24NetworkHeader  h(MasterNode, TypeInvalid);
  send_message(h.from_node, h.to_node, h.type, h.id, 0, 0);
}

void handle_network()
{
  RF24NetworkHeader h;
  int c = network.read(h, buffer, sizeof(buffer));

  send_message(h.from_node, h.to_node, h.type, h.id, buffer, c);
}

void loop(void){
  // 1. Check the network regularly
  network.update();
  
  // 2. Simply forware all message from network to serial line
  if ( network.available() ) {
    handle_network();
  }

  // 3. dispatch message from serial line to network nodes
  if (0 && Serial.available()) {
    if (!handle_serial_line()) {
      // clear all rxd buffer and issuing warings
      discard_serial_line_buffer();
    }
  }
}