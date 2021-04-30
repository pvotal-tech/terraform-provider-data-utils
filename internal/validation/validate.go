package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// InEnum valid the provided value belong to the permitted enumeration
func InEnum(enum []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		diagnostics := diag.Diagnostics{}
		v, ok := i.(string)
		if !ok {
			diagnostics = append(diagnostics, diag.Diagnostic{
				Severity:      diag.Error,
				AttributePath: path,
				Summary:       "Invalid value provided",
				Detail:        fmt.Sprintf("expected type of %s to be string", reflect.TypeOf(i).String()),
			})
			return diagnostics
		}

		for _, str := range enum {
			if v == str || (ignoreCase && strings.ToLower(v) == strings.ToLower(str)) {
				return diagnostics
			}
		}

		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity:      diag.Error,
			AttributePath: path,
			Summary:       "Invalid value provided",
			Detail:        fmt.Sprintf("expected value %s to be one of %s", v, enum),
		})
		return diagnostics
	}
}
