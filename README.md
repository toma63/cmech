# cmech
A simple celestial mechanics package in golang

## Body
A struct which models a celestial body in the solar system.
Mass in kg, state vector in cartesian coordinates (km-s).

## getsv
Uses the JPL Horizons telnet server to get a mass and state vector for a celestial body.  Includes a command line demo application.

