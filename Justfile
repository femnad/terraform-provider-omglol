build:
    #!/usr/bin/env bash
    set -euEo pipefail
    go build -o terraform-provider-omglol
    mv terraform-provider-omglol ~/.terraform.d/plugins/registry.terraform.io/femnad/omglol/0.1.0/linux_amd64/
    pushd contrib
    if [ -f .terraform.lock.hcl ]
    then
        rm .terraform.lock.hcl
    fi
    terraform init
    popd
