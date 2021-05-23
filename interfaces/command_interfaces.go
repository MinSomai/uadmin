package interfaces

type ICommand interface {
	Proceed(subaction string, args []string)
	GetHelpText() string
}
