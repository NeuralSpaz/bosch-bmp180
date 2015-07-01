//Texas Instruments INA219 high side current sensor

package bmp180

import (
	"fmt"

	"github.com/NeuralSpaz/i2c"
)

const (
	BMP180_ADDRESS         = 0x77
	BMP180_CALREG_AC1      = 0xAA
	BMP180_CALREG_AC2      = 0xAC
	BMP180_CALREG_AC3      = 0xAE
	BMP180_CALREG_AC4      = 0xB0
	BMP180_CALREG_AC5      = 0xB2
	BMP180_CALREG_AC6      = 0xB4
	BMP180_CALREG_B1       = 0xB6
	BMP180_CALREG_B2       = 0xB8
	BMP180_CALREG_MB       = 0xBA
	BMP180_CALREG_MC       = 0xBC
	BMP180_CALREG_MD       = 0xBE
	BMP180_CONTROL         = 0xF4
	BMP180_TEMP            = 0xF6
	BMP180_PRESSURE        = 0xF6
	BMP180_READTEMPCMD     = 0x2E
	BMP180_READPRESSURECMD = 0x34
)

type BMP180 struct {
	Dev     i2c.I2CBus
	Init    bool
	Address uint8

	oss uint

	ac1, ac2, ac3      int16
	ac4, ac5, ac6      uint16
	b1, b2, mb, mc, md int16
	b5                 int32

	Temp     float64
	Pressure float64
	Altitude float64
}

func (d *BMP180) String() string {
	return fmt.Sprintf("Temprature %f", d.Temp)
}

func New(deviceAdress uint8, i2cbus byte) *BMP180 {
	deviceBus := i2c.NewI2CBus(i2cbus)
	d := &BMP180{
		Dev:     deviceBus,
		Address: deviceAdress,
	}
	return d
}

// Fetch all values from BMP180
func Fetch(d *BMP180) error {

	if !d.Init {
		err := calibration(d)
		if err != nil {
			return err
		}
	}

	// Measure Temprature
	if err := d.Dev.WriteByteToReg(d.Address, BMP180_TEMP, BMP180_READTEMPCMD); err != nil {
		return err
	}

	// Get TEMP
	ut, err := d.Dev.ReadWordFromReg(d.Address, BMP180_TEMP)
	if err != nil {
		return err
	}
	UT := int32(ut)

	X1 := ((UT - int32(d.ac6)) * int32(d.ac5)) >> 15
	X2 := (int32(d.mc) * 2048) / (X1 + int32(d.md))
	B5 := X1 + X2
	CT := (B5 + 8) >> 4
	d.Temp = float64(CT) * 0.1
	return nil
}

func calibration(d *BMP180) error {

	ac1, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC1)
	if err != nil {
		return err
	}
	d.ac1 = int16(ac1)

	ac2, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC2)
	if err != nil {
		return err
	}
	d.ac2 = int16(ac2)

	ac3, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC3)
	if err != nil {
		return err
	}
	d.ac3 = int16(ac3)

	ac4, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC4)
	if err != nil {
		return err
	}
	d.ac4 = uint16(ac4)

	ac5, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC5)
	if err != nil {
		return err
	}
	d.ac5 = uint16(ac5)

	ac6, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_AC6)
	if err != nil {
		return err
	}
	d.ac6 = uint16(ac6)

	b1, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_B1)
	if err != nil {
		return err
	}
	d.b1 = int16(b1)

	b2, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_B2)
	if err != nil {
		return err
	}
	d.b2 = int16(b2)

	mb, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_MB)
	if err != nil {
		return err
	}
	d.mb = int16(mb)

	mc, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_MC)
	if err != nil {
		return err
	}
	d.mc = int16(mc)

	md, err := d.Dev.ReadWordFromReg(d.Address, BMP180_CALREG_MD)
	if err != nil {
		return err
	}
	d.md = int16(md)

	d.Init = true
	return nil
}
