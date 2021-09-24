terraform {
  required_providers {
    data-utils = {
      version = "0.0.1"
      source  = "pvotal-tech/data-utils"
    }
  }
}

provider "data-utils" {

}

locals {
  yamlTestCases = {
    for file in fileset("${path.module}/../resources/success", "*.yaml") :
    file => yamldecode(file("${path.module}/../resources/success/${file}"))
  }
  jsonTestCases = {
    for file in fileset("${path.module}/../resources/success", "*.json") :
    file => jsondecode(file("${path.module}/../resources/success/${file}"))
  }


}

data "deep_merge" "yaml" {
  for_each = local.yamlTestCases
  provider = data-utils
  inputs = [for input in each.value.inputs: yamlencode(input)]
  config {
    format                          = "YAML"
    with_override                   = each.value.config.with_override
    with_append_slice               = each.value.config.with_append_slice
    with_overwrite_with_empty_value = each.value.config.with_overwrite_with_empty_value
    with_slice_deep_copy            = each.value.config.with_slice_deep_copy
  }
}

data "deep_merge" "json" {
  for_each = local.jsonTestCases
  provider = data-utils
  inputs = [for input in each.value.inputs: jsonencode(input)]
  config {
    format                          = "JSON"
    with_override                   = each.value.config.with_override
    with_append_slice               = each.value.config.with_append_slice
    with_overwrite_with_empty_value = each.value.config.with_overwrite_with_empty_value
    with_slice_deep_copy            = each.value.config.with_slice_deep_copy
  }
}

output "yaml_results" {
  value = {for yaml in keys(local.yamlTestCases): yaml => yamldecode(data.deep_merge.yaml[yaml].output)}
}


output "json_results" {
  value = {for json in keys(local.jsonTestCases): json => jsondecode(data.deep_merge.json[json].output)}
}
