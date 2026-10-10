package controllers

import (
	"reflect"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type FiberControllerBeamProvider struct {
	FiberController
}

var _ BeamProvider = &FiberControllerBeamProvider{}

var _, _ = any(&FiberControllerBeamProvider{}).(BeamProvider)

// ProvideBeam implements [core.BeamProvider].
func (this *FiberControllerBeamProvider) ProvideBeam(instance reflect.Type, container InitializingContext) (parent Beam, controller []Beam) {
	//nolint:errcheck //  TODO: Should log the error
	Try(func() {
		parentBuilder := RaiseVal(NewBeamBuilder(instance))
		parent = parentBuilder.Build()
		fiberControllerBuilder := RaiseVal(NewBeamBuilder(FiberControllerConfig{}.New(parent, container)))
		controller = []Beam{fiberControllerBuilder.Build()}
	})
	return parent, controller
}
