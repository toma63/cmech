package main

import ("fmt"
	"github.com/toma63/cmech/getsv"
	"time"
	"flag"
)



// main function
func main() {
	object := flag.String("object", "", "The name of a celestial object")
	flag.Parse()
	now := time.Now()
	body := getsv.GetSV(*object, now)
	fmt.Printf("Mass of %s (kg): %e\n", *object, body.Mass)
	fmt.Printf("Position (km): x: %e y: %e z: %e\n", body.X, body.Y, body.Z)
	fmt.Printf("Velocity (km/s): x: %e y: %e z: %e\n", body.Vx, body.Vy, body.Vz)
}
