package core

type FieldFormOptions struct {
	Name            string
	Initial         interface{}
	DisplayName     string
	Validators      *ValidatorRegistry
	Choices         *FieldChoiceRegistry
	HelpText        string
	WidgetType      string
	ReadOnly        bool
	Required        bool
	WidgetPopulate  func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{}
	IsFk            bool
	Autocomplete    bool
	ListFieldWidget string
}

func (ffo *FieldFormOptions) GetName() string {
	return ffo.Name
}

func (ffo *FieldFormOptions) IsItFk() bool {
	return ffo.IsFk
}

func (ffo *FieldFormOptions) GetListFieldWidget() string {
	return ffo.ListFieldWidget
}

func (ffo *FieldFormOptions) GetIsAutocomplete() bool {
	return ffo.Autocomplete
}

func (ffo *FieldFormOptions) GetWidgetPopulate() func(widget IWidget, renderContext *FormRenderContext, currentField *Field) interface{} {
	return ffo.WidgetPopulate
}

func (ffo *FieldFormOptions) GetInitial() interface{} {
	return ffo.Initial
}

func (ffo *FieldFormOptions) GetDisplayName() string {
	return ffo.DisplayName
}

func (ffo *FieldFormOptions) GetValidators() *ValidatorRegistry {
	if ffo.Validators == nil {
		return NewValidatorRegistry()
	}
	return ffo.Validators
}

func (ffo *FieldFormOptions) GetChoices() *FieldChoiceRegistry {
	return ffo.Choices
}

func (ffo *FieldFormOptions) GetHelpText() string {
	return ffo.HelpText
}

func (ffo *FieldFormOptions) GetWidgetType() string {
	return ffo.WidgetType
}

func (ffo *FieldFormOptions) GetReadOnly() bool {
	return ffo.ReadOnly
}

func (ffo *FieldFormOptions) GetIsRequired() bool {
	return ffo.Required
}

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
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "TextareaFieldOptions",
		WidgetType: "textarea",
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
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "DynamicTypeFieldOptions",
		WidgetType: "dynamic",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "ForeignKeyFieldOptions",
		WidgetType: "foreignkey",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:            "ForeignKeyWithAutocompleteFieldOptions",
		WidgetType:      "foreignkey",
		Autocomplete:    true,
		ListFieldWidget: "fklink",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:            "ForeignKeyReadonlyFieldOptions",
		WidgetType:      "foreignkey",
		ReadOnly:        true,
		Autocomplete:    true,
		ListFieldWidget: "fklink",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "TextareaReadonlyFieldOptions",
		WidgetType: "textarea",
		ReadOnly:   true,
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "EmailFieldOptions",
		WidgetType: "email",
	})
	UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "URLFieldOptions",
		WidgetType: "url",
	})
}
