# gover

A web-controlled rover running on a Raspberry Pi.

### Parts used

General:
  - A Raspberry Pi (Recommended: Zero W)
  - Sensor: HC-SR04 Ultrasonic Sensor Distance Module
  - Motor Controller Board: L298N Dual H Bridge 
  - A power bank as supply for the pi (with an micro usb output)
  - A battery pack as supply for the motor board

Things you might find together in one set:
  - about 20 small jumper wires (this rover needs mostly male to female cables, but others are also handy to have)
  - A breadboard for the sensor and the LED (you could also solder it)
  - 1k Ohm resistors (one for the sensor echo, one for the front LED)
  - A LED

The rest:
  - Chassis
  - 4 DC Motors (if they come without cables, you can solder some male to male cables to them)

For the last part I would recomment [this set](https://www.amazon.com/Diymore-Wheels-Chassis-Encoder-Arduino/dp/B01N3PCWHC/) <sup>[link for germany](https://www.amazon.de/diymore-Roboter-Geschwindigkeit-Batterie-Raspberry/dp/B07JK33HVL/)</sup>.


### Settings

There is no settings file, some options need to be set in the code.

To get an overview of which pins are available, run `gpio readall` on your Raspberry Pi. Note the difference between BCM and physical numbering. The values in the code should be the BCM ones.
The physical values are for the Raspberry Pi Zero.

|Option Name   |Value           |Description                                                                                         |
|--------------|----------------|----------------------------------------------------------------------------------------------------|
|Front LED     |17 (physical 11)|The BCM number for the front LED. The other connection of the LED must be connected to a Ground Pin.|
|Sensor Trigger|23 (physical 16)|BCM number for the TRIG output                                                                      |
|Sensor Echo   |24 (physical 18)|BCM number for the ECHO input                                                                       |
|Motor IN1     |5  (physical 29)|BCM number for the IN1 connection of the motor controller board |
|Motor IN2     |6  (physical 31)|BCM number for the IN2 connection of the motor controller board                                                                                                     |
|Motor IN3     |16 (physical 36)|BCM number for the IN3 connection of the motor controller board                                     |
|Motor IN4     |26 (physical 37)|BCM number for the IN4 connection of the motor controller board                                                                                                     |


