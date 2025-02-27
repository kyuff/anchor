// Package env is used to set environment variables, primarily during a test run.
// It uses two concepts to do that.
//
// *Actions*
// Express how the value is set.
// Can either be:
// - **Override** the current value, no matter if it was there before.
// - **Default** is used if there where no value already.
//
// *Input Types*
// How is new values read.
// Can be one of:
// - **File** path to the file where data is stored ina KEY=value multi-lineformat
// - **EnvKeyFile** Environment key that holds a file path to a file with a values in
// - **KeyValue** directly pass key and value in as arguments.
package env
