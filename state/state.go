package state

import "lea/key"

const CHUNKSIZE = 4 // Bytes

/*    FLAGS     */

var CYPHERMODE string = "ecb"

var FILEPATH string = ""

var RECURSION bool = false
var VERBOSE bool = false

/*
--------------
    STATES
--------------
*/

// Active key and seed
var Key key.KeyMaterial

// Wipe iterations
var Iterations int

var Mode string = "Decryption"
