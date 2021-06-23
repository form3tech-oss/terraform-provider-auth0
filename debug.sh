#!/bin/bash

PLUGINS=~/.terraform.d/plugins
VERSION=0.6.2

make build-only && \
  mv $PLUGINS/terraform-provider-auth0 $PLUGINS/form3.tech/providers/auth0/$VERSION/linux_amd64 && \
  $GOBIN/dlv exec --headless --listen=:2345 --api-version=2 \
  $PLUGINS/form3.tech/providers/auth0/$VERSION/linux_amd64/terraform-provider-auth0_v$VERSION -- --debug