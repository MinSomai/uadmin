package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"
	"strings"
)

type CommandRegistry struct {
	actions map[string]interfaces.CommandInterface
}

func (r CommandRegistry) addAction(name string, command interfaces.CommandInterface) {
	r.actions[name] = command
}

func (r CommandRegistry) isRegisteredCommand(name string) bool {
	_, err := r.actions[name]
	return !!err
}

func (r CommandRegistry) runAction(name string) {
	action, _ := r.actions[name]
	action.ParseArgs()
	action.Proceed()
}

func (r CommandRegistry) MakeHelpText() string{
	var helpParts []string
	var i int = 1
	for action, handler := range r.actions {
		helpParts = append(helpParts, fmt.Sprintf("%d. %s - %s", i, action, handler.GetHelpText()))
		i += 1
	}
	return strings.Join(helpParts, "\n")
}