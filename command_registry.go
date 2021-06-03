package uadmin

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"
	"strings"
)

type CommandRegistry struct {
	Actions map[string]interfaces.ICommand
}

func (r CommandRegistry) addAction(name string, command interfaces.ICommand) {
	r.Actions[name] = command
}

func (r CommandRegistry) isRegisteredCommand(name string) bool {
	_, err := r.Actions[name]
	return !!err
}

func (r CommandRegistry) runAction(command string, subaction string, args []string) error {
	action, _ := r.Actions[command]
	return action.Proceed(subaction, args)
}

func (r CommandRegistry) MakeHelpText() string{
	var helpParts []string
	var i int = 1
	for action, handler := range r.Actions {
		helpParts = append(helpParts, fmt.Sprintf("%d. %s - %s", i, action, handler.GetHelpText()))
		i += 1
	}
	return strings.Join(helpParts, "\n")
}
