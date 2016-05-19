#include <SPI.h>
#include <RF24.h>
#include <RF24Network.h>

#include "protocol.h"




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
    }

    Serial.flush();
}


inline
void send_hello()
{
  RF24NetworkHeader  h(MasterNode, TypeHello);

  uint8_t f[2] = {0};
  f[0] = '0' + sizeof(Header);
  f[1] = 0;
//  send_message(MasterNode, MasterNode, h.type, h.id, "hello", sizeof("hello"));
   send_message(MasterNode, MasterNode, h.type, h.id, f, 2);
}

void setup(void)
{
    Serial.begin(MasterSerialBauds);
    setup_RF(MasterNode);
    delay(5);
    send_hello();
}
 


inline
void discard_serial_line_buffer()
{
  while (Serial.available()) {
    Serial.read();
  }
  Serial.flush();
}



void handle_serial_line(const Header* h, uint8_t* payload)
{
    bool failed = false;
    if (!network.is_valid_address(h->from_node)) {
	failed = true;
    }
    
    if (!network.is_valid_address(h->to_node)) {
	failed = true;
    }

    if (!CheckTypeValidation((T)h->type)) {
	failed = true;
    }


    if (failed) {
	send_message(MasterNode, MasterNode, TypeInvalid, 0, (uint8_t*)&h, sizeof(Header));

	discard_serial_line_buffer();
	return;
    }

    
    RF24NetworkHeader header(h->to_node, h->type);
    network.write(header, payload, h->payload_size);

    return;
}


void read_serial_line()
{
    static uint8_t SerialBuffer[BufferSize+sizeof(Header)];
    static uint8_t SerialBufferPos = 0;
    int comming = -1;


    if (!Serial.available()) {
	return;
    }

    while (Serial.peek() >= 0) {
	SerialBuffer[SerialBufferPos++] = Serial.read();

	if (SerialBufferPos < sizeof(Header)) {
	    continue;
	}
	if (SerialBufferPos < (SerialBuffer[7] + sizeof(Header))) {
	    continue;
	}
	

	Header* h = (Header*) SerialBuffer;
	uint8_t* payload = SerialBuffer + sizeof(Header);
	send_message(MasterNode, MasterNode, TypeEcho, 0, SerialBuffer, SerialBufferPos);
	SerialBufferPos = 0;
	handle_serial_line(h, payload);
    }
    return;

}


void handle_network()
{
    network.update();
    if (!network.available() ) {
	return;
    }
    
    uint8_t buffer[BufferSize];
    RF24NetworkHeader h;
    int c = network.read(h, buffer, sizeof(buffer));

    send_message(h.from_node, h.to_node, h.type, h.id, buffer, c);
}

void loop(void){
    
    // 2. Simply forware all message from network to serial line

    handle_network();


    // 3. dispatch message from serial line to network nodes
    read_serial_line();
}