package cmech

import ("math"
)

// celestial body
type Body struct {
	X float64
	Y float64
	Z float64
	Vx float64
	Vy float64
	Vz float64
	Mass float64
}

// gravitational constant
const G float64 = 6.674e-11 // N m**2 / kg**2

// compute the Euclidian distance between two bodies
func (body *Body) Dist(other *Body) float64 {
	dx := body.X - other.X
	dy := body.Y - other.Y
	dz := body.Z - other.Z
	return math.Sqrt((dx * dx) + (dy * dy) + (dz * dz))
}

// update state of body and other based on pairwise interaction
// G m1 m2 / r ** 2
// meters, seconds
func (body *Body) Update(other *Body, timestep float64) {

	dx := body.X - other.X
	dy := body.Y - other.Y
	dz := body.Z - other.Z

	r := math.Sqrt((dx * dx) + (dy * dy) + (dz * dz))
	
	ab := -(G * other.Mass) / (r * r) 
	ao := (G * body.Mass) / (r * r) 

	axb := ab * (dx / r) 
	ayb := ab * (dy / r)
	azb := ab * (dz / r)

	axo := ao * (dx / r) 
	ayo := ao * (dy / r)
	azo := ao * (dz / r)

	// update locations based on velocity
	body.X = body.X + (body.Vx / timestep)
	body.Y = body.Y + (body.Vy / timestep)
	body.Z = body.Z + (body.Vz / timestep)

	other.X = other.X + (other.Vx / timestep)
	other.Y = other.Y + (other.Vy / timestep)
	other.Z = other.Z + (other.Vz / timestep)

	// update velocities based on acceleration
	body.Vx = body.Vx + (axb / timestep)
	body.Vy = body.Vy + (ayb / timestep)
	body.Vz = body.Vz + (azb / timestep)
	
	other.Vx = other.Vx + (axo / timestep)
	other.Vy = other.Vy + (ayo / timestep)
	other.Vz = other.Vz + (azo / timestep)
}
