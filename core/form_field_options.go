package core

type UadminFormConfigurableOptionRegistry struct {
	Options map[string]IFieldFormOptions
}

func (c *UadminFormConfigurableOptionRegistry) AddFieldFormOptions(formOptions IFieldFormOptions) {
	c.Options[formOptions.GetName()] = formOptions
}

func (c *UadminFormConfigurableOptionRegistry) GetFieldFormOptions(formOptionsName string) IFieldFormOptions {
	ret, _ := c.Options[formOptionsName]
	return ret
}

var UadminFormCongirurableOptionInstance *UadminFormConfigurableOptionRegistry

func init() {
	UadminFormCongirurableOptionInstance = &UadminFormConfigurableOptionRegistry{
		Options: make(map[string]IFieldFormOptions),
	}
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "RequiredSelectFieldOptions",
		WidgetType: "select",
		Required:   true,
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ContentTypeFieldOptions",
		WidgetType: "contenttypeselector",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "SelectFieldOptions",
		WidgetType: "select",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ReadonlyTextareaFieldOptions",
		WidgetType: "textarea",
		ReadOnly:   true,
	})
	fieldChoiceRegistry := FieldChoiceRegistry{}
	fieldChoiceRegistry.Choices = make([]*FieldChoice, 0)
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:        "UsernameOptions",
		Initial:     "InitialUsername",
		DisplayName: "Username",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ImageFormOptions",
		WidgetType: "image",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "OTPRequiredOptions",
		WidgetType: "hidden",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:     "ReadonlyField",
		ReadOnly: true,
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "PasswordOptions",
		WidgetType: "password",
		HelpText:   "To reset password, clear the field and type a new password.",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ChooseFromSelectOptions",
		WidgetType: "choose_from_select",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DateTimeFieldOptions",
		WidgetType: "datetime",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DatetimeReadonlyFieldOptions",
		WidgetType: "datetime",
		ReadOnly:   true,
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:     "RequiredFieldOptions",
		Required: true,
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkFieldOptions",
		IsFk:       true,
		WidgetType: "fklink",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkReadonlyFieldOptions",
		IsFk:       true,
		ReadOnly:   true,
		WidgetType: "fklink",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "FkRequiredFieldOptions",
		IsFk:       true,
		Required:   true,
		WidgetType: "fklink",
	})
}
