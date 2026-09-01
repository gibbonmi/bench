package example

// Compiled is the package's executed declaration. Go drops a directory whose every file
// a constraint excludes out of a recursive pattern, so without it the stress-tagged file
// beside this one would leave the package unselected instead of unbuilt.
const Compiled = true
