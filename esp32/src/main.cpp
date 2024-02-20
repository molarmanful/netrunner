#include "deps.h"

BfButtonManager man(15, 5);
BfButton btn0(BfButton::ANALOG_BUTTON_ARRAY, 0);
BfButton btn1(BfButton::ANALOG_BUTTON_ARRAY, 1);
BfButton btn2(BfButton::ANALOG_BUTTON_ARRAY, 2);
BfButton btn3(BfButton::ANALOG_BUTTON_ARRAY, 3);
BfButton btn4(BfButton::ANALOG_BUTTON_ARRAY, 4);

void pressH(BfButton *btn, BfButton::press_pattern_t pat) {
  Serial.print(btn->getID());
  Serial.print(" ");
  BfButtonManager::printReading(15);
}

void setup() {
  Serial.begin(9600);

  btn0.onPress(pressH);
  btn1.onPress(pressH);
  btn2.onPress(pressH);
  btn3.onPress(pressH);
  btn4.onPress(pressH);

  man.addButton(&btn0, 300, 550);
  man.addButton(&btn1, 560, 900);
  man.addButton(&btn2, 910, 1200);
  man.addButton(&btn3, 1210, 2000);
  man.addButton(&btn4, 2010, 4096);
}

void loop() { man.loop(); }
