// +build test

package log

import "io/ioutil"

func init() {
	SetOutput(ioutil.Discard)
}
