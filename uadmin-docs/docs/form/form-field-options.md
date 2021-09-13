# Form field options

It simplifies to configure fields for forms. For example, if you want to make some field readonly, just specify it in the gorm model:
```go
type Language struct {
	Code           string `uadminform:"ReadonlyField"`
}
```
The list of predefined field form options is here:
```go
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
```
You can easily add your own field form option:
```go
uadmincore.UadminFormCongirurableOptionInstance.AddFieldFormOptions(&FieldFormOptions{
		Name:       "YOUROWNFIELDFORMOPTIONS",
		WidgetType: "select",
		Required:   true,
	})
```
