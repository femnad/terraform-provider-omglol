build:
    #!/usr/bin/env bash
    version=$(git tag | tail -n 1 | cut -c2-)
    set -euEo pipefail
    go build -o terraform-provider-omglol
    mkdir -p ~/.terraform.d/plugins/registry.terraform.io/femnad/omglol/$version/linux_amd64/
    mv terraform-provider-omglol ~/.terraform.d/plugins/registry.terraform.io/femnad/omglol/$version/linux_amd64/
    pushd contrib
    if [ -f .terraform.lock.hcl ]
    then
        rm .terraform.lock.hcl
    fi
    terraform init
    popd
