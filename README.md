# terraform-provider-data-utils
Data oriented terraform utils provider


## Usage

### Provider setup
```
terraform {
  required_providers {
    data-utils = {
      source = "pvotal-tech/terraform-provider-data-utils"
      version = "~> 0.0.8"
    }
  }
}

provider "data-utils" {
  # Configuration options
}

```

### Using the data

```
data "deep_merge" "merged" {
  inputs = [
    yamlencode(var.object1),
    yamlencode(var.object2),
  ]
  config {
    format = "YAML"
    with_override = true
    with_append_slice = true
    with_overwrite_with_empty_value = true
    with_slice_deep_copy = true
  }
}

output "merged" {
  value = data.deep_merge.merged.output
}
```