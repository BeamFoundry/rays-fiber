package internal

import (
	"reflect"

	. "github.com/BeamFoundry/rays/pkg/core"
)

// TODO: Generalize that even for server Listener and Config ???
type FiberMwConfig[T any] struct {
	ConfigurationProperties `Name:"rays.fiber"`
	servers                 map[string]*T
}

func (this *FiberMwConfig[T]) GetConfig(instance any) *T {
	config := this.servers["primary"]

	if field := FindFieldForInterface[Server](reflect.TypeOf(instance)); field != nil {
		if server := GetFieldFirstTag(field, []string{"_", "Name", "Server"}); server != "" {
			if cfg, ok := this.servers[server]; ok {
				config = cfg
			}
		}
	}
	return config
}
