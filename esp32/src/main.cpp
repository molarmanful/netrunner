#include "deps.h"

RingBuf<uint8_t, 10> r;
int a, b;

void setup() { Serial.begin(9600); }

void loop() {
  r.pushOverwrite(analogRead(15) * 9 / 4095);

  int sum = 0;
  for (int i = 0; i < r.size(); i++) {
    sum += r[i];
  }
  b = sum / r.size();

  if (a != b) {
    a = b;
    Serial.println(b);
  }
}
