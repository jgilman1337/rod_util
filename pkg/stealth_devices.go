package pkg

import (
	"math/rand/v2"

	"github.com/go-rod/rod/lib/devices"
)

//TODO: test this

var devListMobile = []devices.Device{
	devices.IPhone4,
	devices.IPhone5orSE,
	devices.IPhone6or7or8,
	devices.IPhone6or7or8Plus,
	devices.IPhoneX,
	devices.BlackBerryZ30,
	devices.Nexus4,
	devices.Nexus5,
	devices.Nexus5X,
	devices.Nexus6,
	devices.Nexus6P,
	devices.Pixel2,
	devices.Pixel2XL,
	devices.LGOptimusL70,
	devices.NokiaN9,
	devices.NokiaLumia520,
	devices.MicrosoftLumia550,
	devices.MicrosoftLumia950,
	devices.GalaxySIII,
	devices.GalaxyS5,
	devices.JioPhone2,
	devices.KindleFireHDX,
	devices.IPadMini,
	devices.IPad,
	devices.IPadPro,
	devices.BlackberryPlayBook,
	devices.Nexus10,
	devices.Nexus7,
	devices.GalaxyNote3,
	devices.GalaxyNoteII,
	devices.MotoG4,
	devices.SurfaceDuo,
	devices.GalaxyFold,
	devices.GalaxyFold,
}

// Gets a random mobile device for go-rod stealth.
func PickRandMobileDevice() devices.Device {
	idx := rand.N(len(devListMobile))
	return devListMobile[idx]
}

var devListDesktop = []devices.Device{
	devices.LaptopWithTouch,
	devices.LaptopWithHiDPIScreen,
	devices.LaptopWithMDPIScreen,
}

// Gets a random desktop device for go-rod stealth.
func PickRandDesktopDevice() devices.Device {
	idx := rand.N(len(devListDesktop))
	return devListDesktop[idx]
}
