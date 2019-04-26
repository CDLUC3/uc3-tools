package git

// ------------------------------------------------------------
// Flags

var FullSHA = false
var Token = ""

// ------------------------------------------------------------
// Unexported symbols

const tokenNotProvided = "can't access GitHub; --token not set (see https://github.com/settings/tokens)"
var inTest = false
