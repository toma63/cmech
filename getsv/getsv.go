package getsv

import ("github.com/ThomasRooney/gexpect"
	"github.com/toma63/cmech"
	"time"
	"strconv"
	"math"
	"regexp"
	"fmt"
	"os"
)

// map names to JPL Horizons object ids
var Name2id = map[string]string {
	"sun": "10",
	"earth": "399",
	"moon": "301",
	"mars": "499",
	"mercury": "199",
	"venus": "299",
	"jupiter": "599",
	"saturn": "699",
	"uranus": "799",
	"neptune": "899",
	"pluto": "999"}

// get state vector for a named celestial body
// return a cmech.Body
// convert to units of km, sec, kg
// takes horizons_time time.Time
// allowing multiple objects to use the same starting time
func GetSV(oname string, horizons_time time.Time) cmech.Body {

	// connect to the JPL Horizons telnet interface
	child, err := gexpect.Spawn("telnet horizons.jpl.nasa.gov 6775")
	if err != nil {
		panic(err)
	}

	// can be either an integer Horizons object number,
	// of one of the planet name mappings provided
	numre := regexp.MustCompile(`^\d+$`)
	id, ok := Name2id[oname]
	if numre.MatchString(oname) {
		id = oname
	} else {
		if ! ok {
			fmt.Printf("Error, no mapping for object %s\n", oname)
			os.Exit(0)
		}
	}
	
	child.Expect("Horizons>")
	child.SendLine(id)

	// get mass
	mass_sl, _ := child.ExpectRegexFind(`Mass\s*\w*,?\s*\(?10\^(\d+)\s*kg\s*\)?\s*(=|~)\s*(\d+\.\d+)`)
	mass_exp, _ := strconv.ParseFloat(mass_sl[1], 64)
	mass_factor := math.Pow(10.0, mass_exp)
	scaled_mass, _ := strconv.ParseFloat(mass_sl[3], 64)
	mass := scaled_mass * mass_factor
	
	// covert time to horizons format
	tnf := horizons_time.Format("2006-Jan-02 15:04")
	tlf := horizons_time.Add(24 * time.Hour).Format("2006-Jan-02 15:04") // one day later

	// get state vector
	child.Expect("<cr>:")
	child.SendLine("e")
	child.Expect("[o,e,v,?] :")
	child.SendLine("v")
	child.Expect("[ <id>,coord,geo  ] :")
	child.SendLine("@0") // solar system barycenter
	child.Expect("[eclip, frame, body ] :")
	child.SendLine("eclip")
	child.Expect("Starting CT")
	child.SendLine(tnf)
	child.Expect("Ending   CT")
	child.SendLine(tlf)
	child.Expect("Output interval [ex: 10m, 1h, 1d, ? ] :")
	child.SendLine("1d")
	child.Expect("Accept default output [ cr=(y), n, ?] :")
	child.SendLine("n")
	child.Expect("Output reference frame [J2000, B1950] :")
	child.SendLine("J2000")
	child.Expect("Corrections [ 1=NONE, 2=LT, 3=LT+S ]  :")
	child.SendLine("1")
	child.Expect("Output units [1=KM-S, 2=AU-D, 3=KM-D] :")
	child.SendLine("1") // km/s
	child.Expect("Spreadsheet CSV format    [ YES, NO ] :")
	child.SendLine("NO")
	child.Expect(" Label cartesian output    [ YES, NO ] :")
	child.SendLine("NO")
	child.Expect("Select output table type  [ 1-6, ?  ] :")
	child.SendLine("2")
	stv_pos_sl, _ := child.ExpectRegexFind(`\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)`)
	stv_vel_sl, _ := child.ExpectRegexFind(`\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)`)

	x, _ := strconv.ParseFloat(stv_pos_sl[1], 64)
	y, _ := strconv.ParseFloat(stv_pos_sl[2], 64)
	z, _ := strconv.ParseFloat(stv_pos_sl[3], 64)
	vx, _ := strconv.ParseFloat(stv_vel_sl[1], 64)
	vy, _ := strconv.ParseFloat(stv_vel_sl[2], 64)
	vz, _ := strconv.ParseFloat(stv_vel_sl[3], 64)

	body := cmech.Body{x,
		y,
		z,
		vx,
		vy,
		vz,
		mass}

	return body
}


