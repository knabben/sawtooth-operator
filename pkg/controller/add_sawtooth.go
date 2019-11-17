package controller

import (
	"github.com/knabben/sawtooth-operator/pkg/controller/sawtooth"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sawtooth.Add)
}
