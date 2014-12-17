package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const (

	// Event Types
	tButton = 1

	// Button Codes
	bcSelect   = 288
	bcL3       = 289
	bcR3       = 290
	bcStart    = 291
	bcUp       = 292
	bcRight    = 293
	bcDown     = 294
	bcLeft     = 295
	bcL2       = 296 // TODO (adammck): L2 and R2 are analog buttons, but only
	bcR2       = 297 //                 seem to be reported as 0 or 1.
	bcL1       = 298
	bcR1       = 299
	bcTriangle = 300
	bcCircle   = 301
	bcCross    = 302
	bcSquare   = 303
	bcPS       = 304

	// Analog Stick Codes
	LeftStickX  = 0
	LeftStickY  = 1
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
}

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/input.h#L24
type inputEvent struct {
	Time  timeVal
	Type  uint16
	Code  uint16
	Value int32
}

// These are stored as int32, so we can stick the Values straight into them
// without casting. They should probably be cast on the way out.
type AnalogStick struct {
	X int32
	Y int32
}

func (as *AnalogStick) String() string {
	return fmt.Sprintf("%+04d, %+04d", as.X, as.Y)
}

type SA struct {
	r          io.Reader

	// Buttons
	Select   bool
	L3       bool
	R3       bool
	Start    bool
	Up       bool
	Right    bool
	Down     bool
	Left     bool
	L2       bool
	R2       bool
	L1       bool
	R1       bool
	Triangle bool
	Circle   bool
	Cross    bool
	Square   bool
	PS       bool

	// Sticks
	LeftStick  *AnalogStick
	RightStick *AnalogStick
}

func New(reader io.Reader) *SA {
	return &SA{
		r:          reader,
		LeftStick:  &AnalogStick{0, 0},
		RightStick: &AnalogStick{0, 0},
	}
}

func (sa *SA) String() string {

	// dpad
	dpad := make([]string, 0, 4)
	if sa.Up {
		dpad = append(dpad, "up")
	}
	if sa.Down {
		dpad = append(dpad, "down")
	}
	if sa.Left {
		dpad = append(dpad, "left")
	}
	if sa.Right {
		dpad = append(dpad, "right")
	}

	// other buttons
	buttons := make([]string, 0, 13)
	if sa.Select {
		buttons = append(buttons, "select")
	}
	if sa.L3 {
		buttons = append(buttons, "L3")
	}
	if sa.R3 {
		buttons = append(buttons, "R3")
	}
	if sa.Start {
		buttons = append(buttons, "start")
	}
	if sa.L2 {
		buttons = append(buttons, "L2")
	}
	if sa.R2 {
		buttons = append(buttons, "R2")
	}
	if sa.L1 {
		buttons = append(buttons, "L1")
	}
	if sa.R1 {
		buttons = append(buttons, "R1")
	}
	if sa.Triangle {
		buttons = append(buttons, "triangle")
	}
	if sa.Circle {
		buttons = append(buttons, "circle")
	}
	if sa.Cross {
		buttons = append(buttons, "cross")
	}
	if sa.Square {
		buttons = append(buttons, "square")
	}
	if sa.PS {
		buttons = append(buttons, "PS")
	}

	return fmt.Sprintf(
		"&Sixaxis{L[%s] R[%s] dpad[%s] buttons[%s]}",
		sa.LeftStick,
		sa.RightStick,
		strings.Join(dpad, ", "),
		strings.Join(buttons, ", "))
}

func (sa *SA) Update(event *inputEvent) {
	switch event.Type {
	case 0:
		// do nothing?

	case tButton:
		v := buttonToBool(event.Value)

		switch event.Code {
		case bcSelect:
			sa.Select = v

		case bcL3:
			sa.L3 = v

		case bcR3:
			sa.R3 = v

		case bcStart:
			sa.Start = v

		case bcUp:
			sa.Up = v

		case bcRight:
			sa.Right = v

		case bcDown:
			sa.Down = v

		case bcLeft:
			sa.Left = v

		case bcL2:
			sa.L2 = v

		case bcR2:
			sa.R2 = v

		case bcL1:
			sa.L1 = v

		case bcR1:
			sa.R1 = v

		case bcTriangle:
			sa.Triangle = v

		case bcCircle:
			sa.Circle = v

		case bcCross:
			sa.Cross = v

		case bcSquare:
			sa.Square = v

		case bcPS:
			sa.PS = v

		default:
			fmt.Println("Unknown button code")
			printEvent(event)
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

		default:
			printEvent(event)
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
