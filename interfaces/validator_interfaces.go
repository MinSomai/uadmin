package interfaces

import "fmt"

type IValidator func(i interface{}, o interface{}) error

type ValidatorRegistry struct {
	Validators map[string]IValidator
}

func (vr *ValidatorRegistry) AddValidator(validatorName string, implementation IValidator) {
	_, exists := vr.Validators[validatorName]
	if exists {
		Trail(WARNING, "You are overriding validator %s", validatorName)
		return
	}
	vr.Validators[validatorName] = implementation
}

func (vr *ValidatorRegistry) GetValidator(validatorName string) (IValidator, error){
	validator, exists := vr.Validators[validatorName]
	if !exists {
		return nil, fmt.Errorf("no %s validator registered", validatorName)
	}
	return validator, nil
}

var UadminValidatorRegistry *ValidatorRegistry

func init() {
	UadminValidatorRegistry = &ValidatorRegistry{
		Validators: make(map[string]IValidator),
	}
}