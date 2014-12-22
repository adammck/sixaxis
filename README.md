# Sixaxis

This is a Go interface to Sony's [Sixaxis] [hardware]. Only Linux is supported
for the time being, because it relies on the Input Subsystem. Not everything is
working yet, but so far it supports:

* [x] Bluetooth (via [sixa] [qtsixa])
* [x] Digial Buttons
* [x] Analog Sticks
* [x] Analog Triggers
* [x] Pressure-sensitive Buttons
* [ ] Accelerometer
* [x] Gyroscope
* [ ] Rumble
* [ ] LEDs


## Documentation

The docs can be found at [godoc.org] [docs], as usual.


## Usage

???


## License

[MIT] [license].


## Author

[Adam Mckaig] [adammck] made this.


[hardware]: https://en.wikipedia.org/wiki/Sixaxis
[qtsixa]:   http://qtsixa.sourceforge.net
[license]:  https://github.com/adammck/sixaxis/blob/master/LICENSE
[docs]:     https://godoc.org/github.com/adammck/sixaxis
[adammck]:  http://github.com/adammck
