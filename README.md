# terraform-provider-auth0 [![Build Status](https://travis-ci.org/kevholditch/terraform-provider-auth0.svg?branch=master)](https://travis-ci.org/kevholditch/terraform-provider-auth0)
A terraform provider for Auth0

# Debugging

In order to debug the provider, first make sure the target workspace is running at least terraform v0.14. Execute following steps:
  * run `./debug.sh` script in this repo root and init target workspace (`terraform init`)
  * attach with your debugger to port `2345`, this will also print required gRPC debug info to stdout
  * note the env variable from the above output (you'll need it in the next step)
  * set breakpoints in code and run a plan (`TF_REATTACH_PROVIDERS='{...}' terraform plan`)
  * re-run the above command as many times as needed (provider process keeps running after the plan has finished)

You may also want to edit your delve config to increase string truncation limit in the debugger (set `max-string-len` in `~/.config/dlv/config.yml`)