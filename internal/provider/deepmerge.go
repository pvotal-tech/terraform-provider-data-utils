package provider

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/imdario/mergo"

	"github.com/pvotal-tech/terraform-provider-data-utils/internal/config"
	"gopkg.in/yaml.v3"

	"github.com/pvotal-tech/terraform-provider-data-utils/internal/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeepMergeConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"format": {
				Description:      "Specify the type of input and output that the merger will be dealing with. Allowed values are: \"JSON\" or \"YAML\".",
				Type:             schema.TypeString,
				Required:         true,
				Optional:         false,
				ValidateDiagFunc: validation.InEnum([]string{"JSON", "YAML"}, false),
			},
			"with_override": {
				Description: "Specify whether the merger should override non-empty current-level attributes with non-empty next-level attributes values.",
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
			},
			"with_append_slice": {
				Description: "Specify whether the merger should append slices instead of overwriting it. This implicitly overrides the `with_slice_deep_copy` behaviour.",
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
			},
			"with_overwrite_with_empty_value": {
				Description: "Specify whether the merger should override non empty current-level attributes with empty next-level attributes values. This implicitly enables the `with_override` flag.",
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
			},
			"with_slice_deep_copy": {
				Description: "Specify whether the merger should merge slice element one by one. This implicitly enables the `with_override` flag.",
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func dataSourceDeepMerge() *schema.Resource {
	return &schema.Resource{
		Description: "The `deep_merge` data source accepts a list of JSON or YAML strings as input and deep merges into a single output, based on the configured merge rules.",

		ReadContext: dataSourceDeepMergeRead,

		Schema: map[string]*schema.Schema{
			"inputs": {
				Description: "A list of YAML or JSON strings that will be merged into the `output` attribute. Inputs are processed in order.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"config": {
				Description: "The merger configuration. Refer to https://github.com/imdario/mergo for more details regarding configuration behaviours.",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem:        dataSourceDeepMergeConfig(),
			},

			"output": {
				Description: "The deep-merged `output`, rendered in YAML or JSON, according to configuration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDeepMergeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inputs := d.Get("inputs")

	conf := config.New(d.Get("config").([]interface{}))

	if conf == nil {
		return diag.FromErr(fmt.Errorf("config block must be specified"))
	}

	output, diagnostics := merge(inputs.([]interface{}), conf)

	if len(diagnostics) > 0 {
		return diagnostics
	}

	err := d.Set("output", output)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(makeId(output))

	return nil
}

func merge(inputs []interface{}, conf *config.Config) (string, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	merged := make(map[string]interface{})
	var unmarshalFunc func([]byte, interface{}) error
	var marshalFunc func(interface{}) ([]byte, error)

	switch conf.Format {
	case "JSON":
		unmarshalFunc = json.Unmarshal
		marshalFunc = json.Marshal
	case "YAML":
		unmarshalFunc = yaml.Unmarshal
		marshalFunc = yaml.Marshal
	}

	for _, strSrc := range inputs {
		src := make(map[string]interface{})

		err := unmarshalFunc([]byte(strSrc.(string)), &src)
		if err != nil {
			diagnostics = append(diagnostics, diag.FromErr(err)...)
			continue
		}
		err = mergo.Merge(&merged, &src, conf.ToMergoConfig()...)
		if err != nil {
			diagnostics = append(diagnostics, diag.FromErr(err)...)
			continue
		}
	}

	if len(diagnostics) > 0 {
		return "", diagnostics
	}
	b, err := marshalFunc(merged)
	if err != nil {
		return "", diag.FromErr(err)
	}
	return string(b), diagnostics

}

func makeId(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
