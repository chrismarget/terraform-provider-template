## terraform-provider-template

This is a super basic terraform provider which writes a single string attribute
to a local temp file. It's called "template" because I plan to use it as a
starting point for other example providers in bug reports, etc...

#### Example Configuration:
```terraform
terraform {
  required_providers {
    template = {
      source = "registry.terraform.io/chrismarget/template"
    }
  }
}

resource "template_a" "example" {
  string_attr = "foo"
}
```

#### Gotchas:
This provider isn't published on the regitry, so using it isn't completely
straightforward. To test it:
- Don't run `terraform init`. You'll just jump straight to `terraform apply`
- Install the provider binary to your local `GOBIN` with `go install`
- Add a configuration like the following to your `~/.terraformrc` so that terraform can find it:
  ```terraform
  provider_installation {
    dev_overrides {
      "chrismarget/template" = "/Users/chrismarget/go/bin"
    }
    direct {}
  }
  ```