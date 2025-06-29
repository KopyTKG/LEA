package state

import (
	"lea/key"
)

const CHUNKSIZE = 4 // Bytes

/*    FLAGS     */

var CYPHERMODE string = "ecb"

var FILEPATH string = ""

var RECURSION bool = false
var VERBOSE bool = false
var ENCRYPT bool = true

/*
--------------
    STATES
--------------
*/

// Active key and seed
var Key key.KeyMaterial

// Wipe iterations
var Iterations int = 3

/*
------------------------
    ERROR HANDLING
------------------------
*/

const (
	Key128 = 128
	Key192 = 192
	Key256 = 256
)

var ValidKeys = map[int]bool{
	Key128: true,
	Key192: true,
	Key256: true,
}
