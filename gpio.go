package gpio

/*
#cgo LDFLAGS: -lrt

#include "gpio.h"

uint16_t DHT11(unsigned pin) {
	uint8_t pulses[5] = {0};
	gpioSetMode(pin, PI_OUTPUT);
	gpioTrigger(pin, 18000, 0);
	gpioSetMode(pin, PI_INPUT);
	usleep(20);
	if (gpioReadPulse(pin, 100, 0) > 100 ||
		gpioReadPulse(pin, 100, 1) > 100) {
		return 0;
	}
	if (gpioReadPulses(pin, 1000, sizeof(pulses) << 3, pulses)) {
		if (pulses[0] + pulses[1] + pulses[2] + pulses[3] == pulses[4]) {
			return (pulses[0] << 8) + pulses[2];
		}
	}
	return 0;
}
*/
import "C"
import "time"

/*
RPi A+ V1.1 GPIO Layout

             3V3 . . 5V
(I2C SDA) GPIO 2 . . 5V
(I2C SCL) GPIO 3 . . GND
          GPIO 4 . . GPIO14 (UART_TXD)
             GND . . GPIO15 (UART_RXD)
          GPIO17 . . GPIO18 (PWM0, PCM_CLK)
          GPIO27 . . GND
          GPIO22 . . GPIO23
             3V3 . . GPIO24
   (MOSI) GPIO10 . . GND
   (MISO) GPIO 9 . . GPIO25
   (SCLK) GPIO11 . . GPIO 8 (CE0_N)
             GND . . GPIO 7 (CE1_N)
I2C ID EEPROM SD . . I2C ID EEPROM SC
          GPIO 5 . . GND
          GPIO 6 . . GPIO12 (PWM1)
   (PWM1) GPIO13 . . GND
   (PWM0) GPIO19 . . GPIO16
          GPIO26 . . GPIO20
             GND . . GPIO21

             ... USB ...
*/

func init() {
	C.gpioInitialise()
}

// Pin type to operate single GPIO pin state, mode and value
type Pin byte

// Mode represents pin mode (see options below)
type Mode byte

// Pin operating mode
const (
	// INPUT (available for read) mode
	INPUT  Mode = C.PI_INPUT
	OUTPUT      = C.PI_OUTPUT
	ALT0        = C.PI_ALT0
	ALT1        = C.PI_ALT1
	ALT2        = C.PI_ALT2
	ALT3        = C.PI_ALT3
	ALT4        = C.PI_ALT4
	ALT5        = C.PI_ALT5
)

// PullState is a pin pull-up/down state
type PullState byte

// Pull states
const (
	// Pull off
	OFF  PullState = C.PI_PUD_OFF
	DOWN           = C.PI_PUD_DOWN
	UP             = C.PI_PUD_UP
)

// SetMode sets pin operating mode
func (pin Pin) SetMode(mode Mode) {
	C.gpioSetMode(C.uint(pin), C.uint(mode))
}

// GetMode gets pin operating mode
func (pin Pin) GetMode() Mode {
	return Mode(C.gpioGetMode(C.uint(pin)))
}

func (pin Pin) Read() bool {
	return C.gpioRead(C.uint(pin)) > 0
}

func (pin Pin) Write(value bool) {
	var intValue C.uint
	if value {
		intValue = 1
	}
	C.gpioWrite(C.uint(pin), intValue)
}

// Trigger the state of the pin to value for t
func (pin Pin) Trigger(t time.Duration, value bool) {
	var intValue C.uint
	if value {
		intValue = 1
	}
	C.gpioTrigger(C.uint(pin), C.uint(t.Nanoseconds()/1000), intValue)
}

// DHT11 reads humidity and temperature from the sensor
func (pin Pin) DHT11() (byte, byte) {
	humidityAndTemperature := int(C.DHT11(C.uint(pin)))
	return byte(humidityAndTemperature >> 8),
		byte(humidityAndTemperature & 255)
}

/*
void gpioSetPullUpDown(unsigned gpio, unsigned pud)
*/
