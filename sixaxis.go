package main

import (
	"fmt"
	"io"
	"os"
	"encoding/binary"
)

const (

	// Event Types
	tButton = 1

	// Button Codes
	bcX = 302

	// Analog Stick Codes
	LeftStickX = 0
	LeftStickY = 1
	RightStickX = 2
	RightStickY = 3

	// Gyro event codes
	GyroX = 4 // tilt left/right
	GyroY = 5 // tilt forwards/backwards

	// Accelerometer codes
	AccelY = 6 // up/down
	// 7??
)

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/time.h#L15
// https://github.com/torvalds/linux/blob/master/include/uapi/asm-generic/posix_types.h#L88
// TODO: Are these also int32 on 64bit systems?
type timeVal struct {
	Sec  int32 // seconds
	Usec int32 // microseconds
};

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/input.h#L24
type inputEvent struct {
	Time  timeVal;
	Type  uint16;
	Code  uint16;
	Value int32;
};

// These are stored as int32, so we can stick the Values straight into them
// without casting. They should probably be cast on the way out.
type AnalogStick struct {
	X int32
	Y int32
}

func (as *AnalogStick) String() string {
	return fmt.Sprintf("(%+04d, %+04d)", as.X, as.Y)
}

type SA struct {
	r io.Reader
	X bool
	LeftStick *AnalogStick
	RightStick *AnalogStick
}

func New(reader io.Reader) *SA {
	return &SA{
		r: reader,
		X: false,
		LeftStick: &AnalogStick{0, 0},
		RightStick: &AnalogStick{0, 0},
	}
}

func (sa *SA) String() string {
	return fmt.Sprintf("&Sixaxis{L%s, R%s}", sa.LeftStick, sa.RightStick)
}

func (sa *SA) Update(event *inputEvent) {
	switch event.Type {
	case 0:
		// do nothing?

	case tButton:
		switch event.Code {
		case bcX:
			sa.X = buttonToBool(event.Value)
			fmt.Println(sa)
		// default:
		// 	fmt.Println("Unknown button code")
		// 	printEvent(event)
		// 	os.Exit(1)
		}

	// gyro? analog?
	case 3:
		switch event.Code {
		case LeftStickX:
			sa.LeftStick.X = event.Value

		case LeftStickY:
			sa.LeftStick.Y = event.Value

		case RightStickX:
			sa.RightStick.X = event.Value

		case RightStickY:
			sa.RightStick.Y = event.Value

		//default:
		//	printEvent(event)
		}

	default:
		fmt.Println("Unknown Event.Type")
		printEvent(event)
	}
}

func buttonToBool(value int32) bool {
	return value == 1
}

func printEvent(event *inputEvent) {
	fmt.Printf("type=%04d, code=%04d, value=%08d\n", event.Type, event.Code, event.Value)
}

func main() {

	// open the device
	f, err := os.Open("/dev/input/event0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var event inputEvent
	sa := New(f)

	// Read events forever
	for {
		binary.Read(sa.r, binary.LittleEndian, &event)
		sa.Update(&event)
		fmt.Println(sa)
	}
}
