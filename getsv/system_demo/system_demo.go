package main

import ("fmt"
	"github.com/toma63/cmech/getsv"
	"time"
	"flag"
	"regexp"
	"strings"
)



// main function
func main() {
	numre := regexp.MustCompile(`^\d+$`)
	objects := flag.String("objects", "", "Colon separated list of names of solar system objects")
	flag.Parse()
	onames := strings.Split(*objects, ":")
	now := time.Now()
	system := getsv.GetCMSystem(onames, now)

	for _, body := range system {
		horznum := ""
		if numre.MatchString(body.Name) {
			horznum = "JPL Horizons object "
		}
		fmt.Printf("Mass of %s%s (kg): %e\n", horznum, body.Name, body.Mass)
		fmt.Printf("Position (km): x: %e y: %e z: %e\n", body.X, body.Y, body.Z)
		fmt.Printf("Velocity (km/s): x: %e y: %e z: %e\n", body.Vx, body.Vy, body.Vz)
	}
}
