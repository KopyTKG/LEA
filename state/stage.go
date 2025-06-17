package state

const CHUNKSIZE = 4 // Bytes

/*    FLAGS     */

var CYPHERMODE string = "ecb"
var KEYLENGTH int = 256

var FILEPATH string = ""

var KEYPATH string = ""
var SEEDPATH string = ""

var RECURSION bool = false
var VERBOSE bool = false

/*
--------------
    STATES
--------------
*/

var ByteKEY []byte
var ByteSEED []byte
var Mode string = "Decryption"
