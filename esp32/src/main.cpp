#include "deps.h"

// RingBuf<uint8_t, 3> r;
uint8_t a, b;

void setup() { Serial.begin(9600); }

void loop() {
  b = analogRead(15) * 9 / 4095;
  // r.pushOverwrite(analogRead(15) * 9 / 4095);
  //
  // uint16_t sum = 0;
  // for (int i = 0; i < r.size(); i++) {
  //   sum += r[i];
  // }
  // b = sum / r.size();

  if (a != b) {
    a = b;
    Serial.println(b);
  }
}
