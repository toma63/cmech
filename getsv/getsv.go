package main

import ("fmt"
	"github.com/ThomasRooney/gexpect"
	"flag"
	"time"
)

// get state vector for a named celestial body
func getsv(oname string) (string, string, []string, []string) {

	// map names to JPL Horizons object ids
	Name2id := map[string]string {
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

	// connect to the JPL Horizons telnet interface
	child, err := gexpect.Spawn("telnet horizons.jpl.nasa.gov 6775")
	if err != nil {
		panic(err)
	}

	id := Name2id[oname]

	child.Expect("Horizons>")
	child.SendLine(id)

	// get mass
	mass_sl, _ := child.ExpectRegexFind(`Mass\s*\w*,?\s*\(?(10\^\d+)\s*kg\s*\)?\s*(=|~)\s*(\d+\.\d+)`)
	mass_factor := mass_sl[1]
	mass := mass_sl[3]

	// current time (first 16 chars)
	tn := time.Now()
	tnf := tn.Format("2006-Jan-02 15:04")
	tlf := tn.Add(24 * time.Hour).Format("2006-Jan-02 15:04") // one day later

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
	child.Expect("Accept default output")
	child.SendLine("y")
	stv_pos_sl, _ := child.ExpectRegexFind(`\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)`)
	stv_pos := stv_pos_sl[1:4]
	stv_vel_sl, _ := child.ExpectRegexFind(`\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)\s+(-?\d\.\d+E[+-]\d+)`)
	stv_vel := stv_vel_sl[1:4]

	return mass_factor, mass, stv_pos, stv_vel
}

// main function
func main() {
	object := flag.String("object", "", "The name of a celestial object")
	flag.Parse()
	mass_factor, mass, stv_pos, stv_vel := getsv(*object)
	fmt.Printf("Mass (x %s kg) of %s: %s\n", mass_factor, *object, mass)
	fmt.Printf("Position: x: %s y: %s z: %s\n", stv_pos[0], stv_pos[1], stv_pos[2])
	fmt.Printf("Velocity: x: %s y: %s z: %s\n", stv_vel[0], stv_vel[1], stv_vel[2])
}
