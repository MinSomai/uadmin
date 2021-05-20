package blueprint

import "github.com/uadmin/uadmin/interfaces"

type Registry struct {
	registeredBlueprints map[string]*interfaces.IBlueprintInterface
}

func (r *Registry) Get(name string) *interfaces.IBlueprintInterface {
	blueprint, _ := r.registeredBlueprints[name]
	return blueprint
}

func (r *Registry) Register(registerBlueprint *interfaces.IBlueprintInterface) {
	r.registeredBlueprints[registerBlueprint.Name] = registerBlueprint
}
