#ifndef __PROTOCOL_H__
#define __PROTOCOL_H__

enum T {
  TypeFirst = 1,
  TypeSerialMessage,
  TypeFPanelControl,
  TypeFPanelStatus,
  TypeInvalid,
  TypeHello,
  TypeEnd,
};

inline bool CheckTypeValidation(T t) { return t > TypeFirst && t < TypeEnd; }

const uint16_t MasterNode = 00;


const uint8_t TypeLen = TypeEnd - TypeFirst;


struct Header{
  uint16_t from_node;
  uint16_t to_node;
  uint16_t id;
  uint8_t type;
  uint8_t payload_size;
};

struct FPanelControl {
  uint8_t pin;
  uint8_t duration;
};

enum {
  FPanelStatus = 0,
  FPanelStatusPower = 1 << 6,
};

const uint8_t BufferSize = 120;

#define MasterSerialBauds 57600
#define RFChannel 90
#define RFPALevel  3 //RF24_PA_MAX
#define RFCRCLen 2 //RF24_CRC_16
#define RFDataRate 0 //RF24_1MBPS

#define PIN_RF_CE 7
#define PIN_RF_CSN 8

#define PIN_RF_SCK 13
#define PIN_RF_MOSI 11
#define PIN_RFMISO 12

#define PIN_POWER_SWITCH 3
#define PIN_POWER_LED 2

#endif
