package forms

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/uadmin/uadmin"
	"github.com/uadmin/uadmin/form"
	"github.com/uadmin/uadmin/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type WidgetTestSuite struct {
	uadmin.UadminTestSuite
}

func NewTestForm() *multipart.Form {
	form1 := multipart.Form{
		Value: make(map[string][]string),
	}
	return &form1
}

func (w *WidgetTestSuite) TestTextWidget() {
	textWidget := &form.TextWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	textWidget.SetName("dsadas")
	textWidget.SetValue("dsadas")
	textWidget.SetRequired()
	renderedWidget := textWidget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	err := textWidget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test"}
	err = textWidget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestNumberWidget() {
	widget := &form.NumberWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"test"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"121"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestEmailWidget() {
	widget := &form.EmailWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"test@example.com"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestURLWidget() {
	widget := &form.URLWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"example.com"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestPasswordWidget() {
	widget := &form.PasswordWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "type=\"password\"")
	widget.SetRequired()
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"12345678901234567890"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestHiddenWidget() {
	widget := &form.HiddenWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas<>")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas&lt;&gt;\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasas"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestDateWidget() {
	widget := &form.DateWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("11/01/2021")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "datetimepicker_dsadas")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"11/02/2021"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestDateTimeWidget() {
	widget := &form.DateTimeWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("11/02/2021 10:04")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"11/02/2021 10:04\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"11/02/2021 10:04"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestTimeWidget() {
	widget := &form.TimeWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("15:05")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"15:05\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestTextareaWidget() {
	widget := &form.TextareaWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render()
	assert.Equal(w.T(), renderedWidget, "<textarea name=\"dsadas\" test=\"test1\">dsadas</textarea>")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestCheckboxWidget() {
	widget := &form.CheckboxWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "checked=\"checked\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	widget.ProceedForm(form1)
	assert.True(w.T(), widget.GetOutputValue() == true)
}

func (w *WidgetTestSuite) TestSelectWidget() {
	widget := &form.SelectWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	widget.OptGroups = make(map[string][]*form.SelectOptGroup)
	widget.OptGroups["test"] = make([]*form.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test1",
		Value: "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test2",
		Value: "dsadas",
	})
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"10:04"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"dsadas"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestNullBooleanWidget() {
	widget := &form.NullBooleanWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.OptGroups = make(map[string][]*form.SelectOptGroup)
	widget.OptGroups["test"] = make([]*form.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test1",
		Value: "yes",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test2",
		Value: "no",
	})
	widget.SetName("dsadas")
	widget.SetValue("yes")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas\" data-placeholder=\"Select\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"no"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSelectMultipleWidget() {
	widget := &form.SelectMultipleWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas"})
	widget.OptGroups = make(map[string][]*form.SelectOptGroup)
	widget.OptGroups["test"] = make([]*form.SelectOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test1",
		Value: "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.SelectOptGroup{
		OptLabel: "test2",
		Value: "dsadas",
	})
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestRadioSelectWidget() {
	widget := &form.RadioSelectWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		Id: "test",
		WrapLabel: true,
	}
	widget.SetName("dsadas")
	widget.SetValue("dsadas")
	widget.OptGroups = make(map[string][]*form.RadioOptGroup)
	widget.OptGroups["test"] = make([]*form.RadioOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.RadioOptGroup{
		OptLabel: "test1",
		Value: "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.RadioOptGroup{
		OptLabel: "test2",
		Value: "dsadas",
	})
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "<li>test<ul id=\"test_0\">")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestCheckboxSelectMultipleWidget() {
	widget := &form.CheckboxSelectMultipleWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		Id: "test",
		WrapLabel: true,
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas"})
	widget.OptGroups = make(map[string][]*form.RadioOptGroup)
	widget.OptGroups["test"] = make([]*form.RadioOptGroup, 0)
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.RadioOptGroup{
		OptLabel: "test1",
		Value: "test1",
	})
	widget.OptGroups["test"] = append(widget.OptGroups["test"], &form.RadioOptGroup{
		OptLabel: "test2",
		Value: "dsadas",
	})
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "<ul id=\"test\">\n  \n  \n  \n    <li>test<ul id=\"test_0\">")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadasdasdas"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err != nil)
	form1.Value["dsadas"] = []string{"test1"}
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestFileWidget() {
	widget := &form.FileWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "type=\"file\"")
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	err := writer.SetBoundary("foo")
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	path := os.Getenv("UADMIN_PATH") + "/tests/file_for_uploading.txt"
	file, err := os.Open(path)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("UADMIN_PATH"), "upload-for-tests"), 0755)
	if err != nil {
		assert.True(w.T(), false, "Couldnt create directory for file uploading", err)
		return
	}
	defer file.Close()
	defer os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("UADMIN_PATH"), "upload-for-tests"))
	part, err := writer.CreateFormFile("dsadas", filepath.Base(path))
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = writer.Close()
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	form1, _ := multipart.NewReader(bytes.NewReader(body.Bytes()), "foo").ReadForm(1000000)
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestClearableFileWidget() {
	widget := &form.ClearableFileWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		InitialText: "test",
		Required: true,
		Id: "test",
		ClearCheckboxLabel: "clear file",
		InputText: "upload your image",
	}
	widget.SetName("dsadas")
	renderedWidget := widget.Render()
	assert.Equal(w.T(), renderedWidget, "<p class=\"file-upload\">test: <br>\nupload your image:\n    <input type=\"file\" name=\"dsadas\" test=\"test1\"></p>")
	widget = &form.ClearableFileWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		InitialText: "test",
		Required: true,
		Id: "test",
		ClearCheckboxLabel: "clear file",
		InputText: "upload your image",
		CurrentValue: &form.URLValue{URL: "https://microsoft.com"},
	}
	widget.SetName("dsadas")
	renderedWidget = widget.Render()
	assert.Equal(w.T(), renderedWidget, "\n    <input type=\"file\" name=\"dsadas\" test=\"test1\">")
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	err := writer.SetBoundary("foo")
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	path := os.Getenv("UADMIN_PATH") + "/tests/file_for_uploading.txt"
	file, err := os.Open(path)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("UADMIN_PATH"), "upload-for-tests"), 0755)
	if err != nil {
		assert.True(w.T(), false, "Couldnt create directory for file uploading", err)
		return
	}
	defer file.Close()
	defer os.RemoveAll(fmt.Sprintf("%s/%s", os.Getenv("UADMIN_PATH"), "upload-for-tests"))
	part, err := writer.CreateFormFile("dsadas", filepath.Base(path))
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	err = writer.Close()
	if err != nil {
		assert.True(w.T(), false)
		return
	}
	form1, _ := multipart.NewReader(bytes.NewReader(body.Bytes()), "foo").ReadForm(1000000)
	err = widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestMultipleHiddenInputWidget() {
	widget := &form.MultipleInputHiddenWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
	}
	widget.SetName("dsadas")
	widget.SetValue([]string{"dsadas", "test1"})
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "value=\"dsadas\"")
	form1 := NewTestForm()
	form1.Value["dsadas"] = []string{"dsadas", "test1"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSplitDateTimeWidget() {
	widget := &form.SplitDateTimeWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		DateAttrs: map[string]string{"test": "test1"},
		TimeAttrs: map[string]string{"test": "test1"},
		TimeFormat: "15:04",
		DateFormat: "Mon Jan _2",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_date\"")
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_time\"")
	form1 := NewTestForm()
	form1.Value["dsadas_date"] = []string{"Mon Jan 12"}
	form1.Value["dsadas_time"] = []string{"10:20"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSplitHiddenDateTimeWidget() {
	widget := &form.SplitHiddenDateTimeWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		DateAttrs: map[string]string{"test": "test1"},
		TimeAttrs: map[string]string{"test": "test1"},
		TimeFormat: "15:04",
		DateFormat: "Mon Jan _2",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_date\"")
	assert.Contains(w.T(), renderedWidget, "name=\"dsadas_time\"")
	form1 := NewTestForm()
	form1.Value["dsadas_date"] = []string{"Mon Jan 12"}
	form1.Value["dsadas_time"] = []string{"10:20"}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

func (w *WidgetTestSuite) TestSelectDateWidget() {
	widget := &form.SelectDateWidget{
		Widget: form.Widget{
			Attrs: map[string]string{"test": "test1"},
			BaseFuncMap: template.FuncMap,
		},
		EmptyLabelString: "choose any",
	}
	widget.SetName("dsadas")
	nowTime := time.Now()
	widget.SetValue(&nowTime)
	renderedWidget := widget.Render()
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_month\"")
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_day\"")
	assert.Contains(w.T(), renderedWidget, "<select name=\"dsadas_year\"")
	form1 := NewTestForm()
	form1.Value["dsadas_month"] = []string{"1"}
	form1.Value["dsadas_day"] = []string{"1"}
	form1.Value["dsadas_year"] = []string{strconv.Itoa(time.Now().Year())}
	err := widget.ProceedForm(form1)
	assert.True(w.T(), err == nil)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestWidget(t *testing.T) {
	uadmin.Run(t, new(WidgetTestSuite))
}
