package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"
)

type BlueprintRegistry struct {
	RegisteredBlueprints map[string]interfaces.IBlueprint
}

func (r BlueprintRegistry) Iterate() <-chan interfaces.IBlueprint {
	chnl := make(chan interfaces.IBlueprint)
	go func() {
		for _, blueprint := range r.RegisteredBlueprints {
			chnl <- blueprint
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func (r BlueprintRegistry) GetByName(name string) (interfaces.IBlueprint, error) {
	blueprint, ok := r.RegisteredBlueprints[name]
	var err error
	if !ok {
		err = fmt.Errorf("Couldn't find blueprint with name %s", name)
	}
	return blueprint, err
}

func (r BlueprintRegistry) Register(blueprint interfaces.IBlueprint) {
	r.RegisteredBlueprints[blueprint.GetName()] = blueprint
}
