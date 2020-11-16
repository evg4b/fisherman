package log

import "io/ioutil"

func init() {
	generalOutput = ioutil.Discard
}
