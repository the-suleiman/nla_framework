package templates

import (
	"fmt"
	"log"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/serenize/snaker"
	"github.com/the-suleiman/nla_framework/types"
	"github.com/the-suleiman/nla_framework/utils"
)

func docIsIntegrationProccess(p types.ProjectType, d *types.DocType) {
	if d.IsOdataIntegration() {
		docIsIntegrationOdataProccess(p, d)
	}
}

func docIsIntegrationOdataProccess(p types.ProjectType, d *types.DocType) {
	sourcePath := getCurrentDir() + "/integrations/odata/odataDoc.go"
	// проверяем возможность того, что путь к шаблону был переопределен внутри документа
	if d.TemplatePathOverride != nil {
		if tmpl, ok := d.TemplatePathOverride["odataDoc.go"]; ok {
			if len(tmpl.Source) > 0 {
				sourcePath = tmpl.Source
			}
		}
	}
	docName := d.Name
	odataName := d.Integrations.Odata.Name
	odataFldNames := []string{}
	for _, fld := range d.Flds {
		odataNameFld := getOdataFld(*d, fld).Name
		if len(odataNameFld) > 0 {
			odataFldNames = append(odataFldNames, odataNameFld)
		}
	}
	localFuncMap := template.FuncMap{
		"LocalProjectPath": func() string { return p.Config.LocalProjectPath },
		"DocNameCamel":     func() string { return snaker.SnakeToCamel(docName) },
		"GetOdataName":     func() string { return odataName },
		"GetOdataFldNames": func() []string { return odataFldNames },
		"IsOdataFld": func(fld types.FldType) bool {
			return len(getOdataFld(*d, fld).Name) > 0
		},
		"GetOdataFldName": func(fld types.FldType) string {
			return getOdataFld(*d, fld).Name
		},
		"GetOdataFldType": func(fld types.FldType) string {
			t := getOdataFld(*d, fld).Type
			if len(t) == 0 {
				return fld.GoType()
			}
			return t
		},
		"CastToGoType": func(fld types.FldType) string {
			fName := strcase.ToCamel(fld.Name)
			// если в описании поля указан способ приведения к типу, то используем его
			if len(getOdataFld(*d, fld).CastToGoType) > 0 {
				return getOdataFld(*d, fld).CastToGoType
			}
			switch fld.Type {
			case types.FldTypeText, types.FldTypeString:
				return fmt.Sprintf("res.%[1]s = cast.ToString(odataDoc.%[1]s)", fName)
			case types.FldTypeInt:
				return fmt.Sprintf("res.%[1]s = cast.ToInt(odataDoc.%[1]s)", fName)
			case types.FldTypeInt64:
				return fmt.Sprintf("res.%[1]s = cast.ToInt64(odataDoc.%[1]s)", fName)
			case types.FldTypeDouble:
				return fmt.Sprintf("res.%[1]s = cast.ToFloat64(odataDoc.%[1]s)", fName)
			case types.FldTypeIntArray:
				return fmt.Sprintf(`res.%[1]s = []int{}
				intSlice%[1]s, err := cast.ToIntSliceE(odataDoc.%[1]s)
				if err == nil {
					res.%[1]s = intSlice%[1]s
				}`, fName)
			case types.FldTypeTextArray:
				return fmt.Sprintf(`res.%[1]s = []string{}
				txtSlice%[1]s, err := cast.ToStringSliceE(odataDoc.%[1]s)
				if err == nil {
					res.%[1]s = txtSlice%[1]s
				}`, fName)
			}
			return "`!!! CastToGoType not found for type: " + fld.Type + " fld: " + fld.Name + "`"
		},
	}
	for k, v := range funcMap {
		localFuncMap[k] = v
	}
	t, err := template.New("odataDoc.go").Funcs(localFuncMap).Delims("[[", "]]").ParseFiles(sourcePath)
	utils.CheckErr(err, "odataDoc.go")
	distPath := fmt.Sprintf("%s/odata", p.DistPath)
	d.Templates["webClient_comp_odataDoc.go"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: snaker.SnakeToCamelLower(d.Name) + ".go"}
}

func getOdataFld(d types.DocType, fld types.FldType) types.OdataFld {
	if odataFldInt, ok := fld.IntegrationData["odata"]; ok {
		if odataFld, ok := odataFldInt.(types.OdataFld); ok {
			return odataFld
		}
		log.Fatalf("docIsIntegrationOdataProccess doc: '%s' fld: '%s' not OdataFld", d.Name, fld.Name)
	}
	return types.OdataFld{}
}
