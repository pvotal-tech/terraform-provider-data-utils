# terraform-provider-data-utils
Data oriented terraform utils provider


## Usage

### Provider setup
```
terraform {
  required_providers {
    merger = {
      source = "app.terraform.io/3rein/3rein-common-provider-tf-merger"
      version = "~> 0.0.1"
    }
  }
}

provider "merger" {
  # Configuration options
}

```

### Using the data

```
data "deep_merge" "merged" {
  inputs = [
    var.object1,
    var.object2,
  ]
  with_override = true
  with_append_slice = true
  with_overwrite_with_empty_value = true
  with_slice_deep_copy = true
}

output "merged" {
  value = data.deep_merge.merged.output
}
```