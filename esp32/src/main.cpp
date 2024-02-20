#include "deps.h"

uint16_t reading, avg, sum;

void setup() { Serial.begin(9600); }

void loop() {
  static unsigned int i = 0;
  reading = BfButtonManager::printReading(15);
  if (reading > 100) {
    sum += reading;
    if (i == 4) {
      avg = sum / 5;
      Serial.print("Average Reading: ");
      Serial.println(avg);
      sum = 0;
    }
    i++;
    if (i > 4)
      i = 0;
  } else {
    sum = 0;
    i = 0;
  }
  delay(100);
}
