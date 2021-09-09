package core

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type ISearchFieldInterface interface {
	Search(afo IAdminFilterObjects, searchString string)
}

type SearchField struct {
	Field        *Field
	CustomSearch func(afo IAdminFilterObjects, searchString string)
}

func (sf *SearchField) Search(afo IAdminFilterObjects, searchString string) {
	if sf.CustomSearch != nil {
		sf.CustomSearch(afo, searchString)
	} else {
		operator := IContainsGormOperator{}
		gormOperatorContext := NewGormOperatorContext(afo.GetFullQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, sf.Field, searchString, true)
		afo.SetFullQuerySet(gormOperatorContext.Tx)
		gormOperatorContext = NewGormOperatorContext(afo.GetPaginatedQuerySet(), afo.GetCurrentModel())
		operator.Build(afo.GetUadminDatabase().Adapter, gormOperatorContext, sf.Field, searchString, true)
		afo.SetPaginatedQuerySet(gormOperatorContext.Tx)
	}
}

func NewSearchFieldRegistryFromGormModel(modelI interface{}) *SearchFieldRegistry {
	if modelI == nil {
		return nil
	}
	ret := &SearchFieldRegistry{Fields: make([]*SearchField, 0)}
	uadminDatabase := NewUadminDatabaseWithoutConnection()
	stmt := &gorm.Statement{DB: uadminDatabase.Db}
	stmt.Parse(modelI)
	gormModelV := reflect.Indirect(reflect.ValueOf(modelI))
	for _, field := range stmt.Schema.Fields {
		uadminTag := field.Tag.Get("uadmin")
		if !strings.Contains(uadminTag, "search") && field.Name != "ID" {
			continue
		}
		uadminField := NewUadminFieldFromGormField(gormModelV, field, nil, true)
		searchField := &SearchField{
			Field: uadminField,
		}
		ret.AddField(searchField)
	}
	return ret
}

type SearchFieldRegistry struct {
	Fields []*SearchField
}

func (sfr *SearchFieldRegistry) GetAll() <-chan *SearchField {
	chnl := make(chan *SearchField)
	go func() {
		defer close(chnl)
		for _, field := range sfr.Fields {
			chnl <- field
		}

	}()
	return chnl
}

func (sfr *SearchFieldRegistry) AddField(sf *SearchField) {
	sfr.Fields = append(sfr.Fields, sf)
}
