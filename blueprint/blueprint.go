package blueprint

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"
)

type BlueprintRegistry struct {
	registeredBlueprints map[string]interfaces.IBlueprint
	//Iterate() <-chan IBlueprint
	//GetByName(name string) IBlueprint
}

func (r BlueprintRegistry) Iterate() <-chan interfaces.IBlueprint {
	chnl := make(chan interfaces.IBlueprint)
	go func() {
		for _, blueprint := range r.registeredBlueprints {
			chnl <- blueprint
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}

func (r BlueprintRegistry) GetByName(name string) (interfaces.IBlueprint, error) {
	blueprint, ok := r.registeredBlueprints[name]
	var err error
	if !ok {
		err = fmt.Errorf("Couldn't find blueprint with name %s", name)
	}
	return blueprint, err
}

func (r BlueprintRegistry) Register(blueprint interfaces.IBlueprint) {
	r.registeredBlueprints[blueprint.GetName()] = blueprint
}
