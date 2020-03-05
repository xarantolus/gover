# gover

A web-controlled rover powered by a Raspberry Pi. This repository contains the source code for the server/rover software. It can be 
controlled over a website using a PC (keyboard) or a Smartphone (tilting), but it also supports video console controllers (e.g. an XBox controller).

## Parts used

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
  - Chassis (no specific one)
  - 4 DC Motors (if they come without cables, you can solder some male to male cables to them)

For the last parts I recommend [this set](https://www.amazon.com/Diymore-Wheels-Chassis-Encoder-Arduino/dp/B01N3PCWHC/) <sup>[link for germany](https://www.amazon.de/diymore-Roboter-Geschwindigkeit-Batterie-Raspberry/dp/B07JK33HVL/)</sup>.

## Pictures & Assembly

This is not an exact guide, but rather a general description of what I did to build this robot.

In the beginning, I started out by wiring and programming the ultrasonic sensor. I tested smaller programs that were able to measure the distance in front of the sensors. It has some limitations that can be worked around in the software: very close (<2cm) and far (>500cm) values are not supported by the hardware, so any measurements outside of this range must be discarded. 

![Ultrasonic sensor wiring](img/sensor.jpg?raw=true)

The next part was getting the motors to work. For that, I wired the motors to the controller board, which can be controlled by the Raspberry Pi using its GPIO pins. Depending on which pins are on, selected motors will start running. In this setup, two motors are wired together (see at the bottom), which means that you can turn the left/right side of the cart on or off. This also allows the robot to change its direction. There is no other way to steer.

![Motor controller board](img/motor-controller-board.jpg?raw=true)

Now I connected all parts together and added them to the chassis. I also joined all software components together (in this project) so they can work together.

![Rover front](img/rover-front.jpg?raw=true) ![Rover side](img/rover-side.jpg?raw=true)


After some time I replaced the Raspberry Pi with a [Raspberry Pi Zero W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) as it is smaller and can run the same software without any changes.

![Rover with Raspberry Pi Zero W](img/rover-raspberrypi-zero.jpg?raw=true)

## Software 
The software is a remote interface to the hardware pins that are connected to the motors. This way, it is not just able to control different motors and an LED, but also measure the distance before it.
Depending on its input and state, the robot moves (or not).

The web clients connect using a WebSocket connection and receive live updates of the current state of the rover (direction, last distance measured). They can change the current state by sending `direction` packages (basically a JSON object) over this connection. 

If the new direction is `forward` and the last distance measurement (that was less than two seconds ago) is less than 25cm, the rover will refuse to go forward. It will also check this distance while driving forward to prevent crashing. The value of 25cm was found out by testing the distance required to stop before crashing into a wall.

Another way to control the rover is using an XBox Controller (the ones with an USB connector that has a wireless connection to the controller). Plugging in the connector will be detected by this software, which will then listen to the signals coming from a controller. It is not recommended to use the web interface and the controller at the same time, as they might cancel out signals coming from the other party.
The buttons `A` and `B` can be used to go forward and backwards, the joystick is used to control left/right direction.

### Supported Directions

- Forward (Button `B` / tilt forward / `W` on PC)
-	Reverse (Button `A` / tilt back / `S` on PC)
- Left (joystick left / tilt left / `A` on PC)
-	Right (joystick right / tilt right / `D` on PC)

There's also the pivot mode where the other side of wheels goes into the opposite direction, which makes turning faster. It is only supported on PCs. 
-	Pivot left (`Q` on PC)
-	Pivot right (`E` on PC)

Reversing with direction is also possible
- Reverse left (Holding `A` and joystick left / `Y` on PC)
-	Reverse right (Holding `A` and joystick right / `C` on PC)

## Building
The software is written in [Go](https://golang.org/) and needs to be cross-compiled to run on a Raspberry Pi. The command can be found in `build_raspberrypi.{cmd,sh}`. It only works if you have Go installed, of course.

The resulting executable can be put on a Raspberry Pi. The `static` directory should be in the same directory as the `gover` executable as the webinterface needs these files. You might need to run this executable with higher privileges as it wants to listen on port `80`. 

### Settings

There is no config file, but some options can be customized in the code (`rover/{sensors,motors,leds}.go` for the output GPIO pins).

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

