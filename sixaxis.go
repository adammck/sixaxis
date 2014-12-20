package sixaxis

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const (

	// Event Types
	tDigital = 1
	tAnalog = 3

	// Digital event codes 0 or 1
	bcSelect   = 288
	bcL3       = 289
	bcR3       = 290
	bcStart    = 291
	bcPS       = 304

	// Analog sticks: -128 to +127
	LeftStickX  = 0
	LeftStickY  = 1
	RightStickX = 2
	RightStickY = 3

	// Gyroscope ... ?
	GyroX = 4 // left/right
	GyroY = 5 // forwards/backwards
	GyroZ = 6 // ???

	// Analog buttons: 0 to 255
	aUp    = 8
	aRight = 9
	aDown  = 10
	aLeft  = 11
	aL2 = 12
	aR2 = 13
	aL1 = 14
	aR1 = 15
	aTriangle = 26
	aCircle = 27
	aCross = 28
	aSquare = 29
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

// Also stored as int32
type Orientation struct {
	X int32
	Y int32
	Z int32
}

// TODO: Scale to the range of values
func (o *Orientation) String() string {
	return fmt.Sprintf("x=%+04d, y=%+04d, z=%+04d", o.X, o.Y, o.Z)
}

type SA struct {
	r io.Reader

	// Digital Buttons
	Select   bool
	L3       bool
	R3       bool
	Start    bool
	PS       bool

	// Analog buttons: 0-255
	Up    int32
	Right int32
	Down  int32
	Left  int32
	L2    int32
	R2    int32
	L1    int32
	R1    int32
	Triangle int32
	Circle   int32
	Cross    int32
	Square   int32

	// Sticks
	// TODO: Should these just be LX/LY/RX/RY?
	LeftStick  *AnalogStick
	RightStick *AnalogStick

	// Gyro
	Orientation *Orientation
}

func New(reader io.Reader) *SA {
	return &SA{
		r:           reader,
		LeftStick:   &AnalogStick{},
		RightStick:  &AnalogStick{},
		Orientation: &Orientation{},
	}
}

// String returns the current state of the controller as a string.
func (sa *SA) String() string {
	s := make([]string, 0, 30)

	// sticks
	if sa.LeftStick.X != 0 {
		s = append(s, fmt.Sprintf("LX=%+04d", sa.LeftStick.X))
	}
	if sa.LeftStick.Y != 0 {
		s = append(s, fmt.Sprintf("LY=%+04d", sa.LeftStick.Y))
	}
	if sa.RightStick.X != 0 {
		s = append(s, fmt.Sprintf("RX=%+04d", sa.RightStick.X))
	}
	if sa.RightStick.Y != 0 {
		s = append(s, fmt.Sprintf("RY=%+04d", sa.RightStick.Y))
	}

	// dpad
	if sa.Up > 0 {
		s = append(s, fmt.Sprintf("up=%d", sa.Up))
	}
	if sa.Down > 0 {
		s = append(s, fmt.Sprintf("down=%d", sa.Down))
	}
	if sa.Left > 0 {
		s = append(s, fmt.Sprintf("left=%d", sa.Left))
	}
	if sa.Right >0 {
		s = append(s, fmt.Sprintf("right=%d", sa.Right))
	}

	// other analogs
	if sa.L1 > 0 {
		s = append(s, fmt.Sprintf("L1=%d", sa.L1))
	}
	if sa.L2 > 0 {
		s = append(s, fmt.Sprintf("L2=%d", sa.L2))
	}
	if sa.R1 > 0 {
		s = append(s, fmt.Sprintf("R1=%d", sa.R1))
	}
	if sa.R2 > 0 {
		s = append(s, fmt.Sprintf("R2=%d", sa.R2))
	}
	if sa.Triangle > 0 {
		s = append(s, fmt.Sprintf("T=%d", sa.Triangle))
	}
	if sa.Circle > 0 {
		s = append(s, fmt.Sprintf("C=%d", sa.Circle))
	}
	if sa.Cross > 0 {
		s = append(s, fmt.Sprintf("X=%d", sa.Cross))
	}
	if sa.Square > 0 {
		s = append(s, fmt.Sprintf("S=%d", sa.Square))
	}

	// digital buttons
	if sa.Select {
		s = append(s, "select")
	}
	if sa.L3 {
		s = append(s, "L3")
	}
	if sa.R3 {
		s = append(s, "R3")
	}
	if sa.Start {
		s = append(s, "start")
	}
	if sa.PS {
		s = append(s, "PS")
	}

	return fmt.Sprintf("&Sixaxis{%s}", strings.Join(s, ", "))
}

// Update changes the state of the controller to reflect the changes in an
// input event. This should be called every time an input event is received.
func (sa *SA) Update(event *inputEvent) {
	switch event.Type {
	case 0:
		// Zero events show up all the time, but never contain any codes or
		// values. I'm guessing that they're sent when the driver has nothing
		// useful to say. So we ignore them.

	case tDigital:
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

		case bcPS:
			sa.PS = v

		default:
			// There are a lot of events which we ignore here, because they're
			// digital representations of the analog buttons. I guess the driver
			// provides these for clients which don't support analog? They're no
			// use to us, anyway.
		}

	case tAnalog:
		switch event.Code {
		case LeftStickX:
			sa.LeftStick.X = event.Value

		case LeftStickY:
			sa.LeftStick.Y = event.Value

		case RightStickX:
			sa.RightStick.X = event.Value

		case RightStickY:
			sa.RightStick.Y = event.Value

		case GyroX:
			sa.Orientation.X = event.Value

		case GyroY:
			sa.Orientation.Y = event.Value

		case GyroZ:
			sa.Orientation.Z = event.Value

		case aUp:
			sa.Up = event.Value

		case aRight:
			sa.Right = event.Value

		case aDown:
			sa.Down = event.Value

		case aLeft:
			sa.Left = event.Value

		case aL2:
			sa.L2 = event.Value

		case aR2:
			sa.R2 = event.Value

		case aL1:
			sa.L1 = event.Value

		case aR1:
			sa.R1 = event.Value

		case aTriangle:
			sa.Triangle = event.Value

		case aCircle:
			sa.Circle = event.Value

		case aCross:
			sa.Cross = event.Value

		case aSquare:
			sa.Square = event.Value

		default:
			//fmt.Println("Unknown event code!")
			//printEvent(event)
		}

	default:
		//fmt.Println("Unknown Event.Type")
		//printEvent(event)
	}
}

// Run loops forever, keeping the state of the controller up to date. This
// should be called in a goroutine.
func (sa *SA) Run() {
	var event inputEvent

	for {
		binary.Read(sa.r, binary.LittleEndian, &event)
		sa.Update(&event)
	}
}

func buttonToBool(value int32) bool {
	return value == 1
}


func printEvent(event *inputEvent) {
	fmt.Printf("type=%04d, code=%04d, value=%08d\n", event.Type, event.Code, event.Value)
}
