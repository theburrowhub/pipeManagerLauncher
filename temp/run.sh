#!/bin/bash

set -e

export PIPELINE_VARIABLE_AUTHOR="Sergio Tejón"
export PIPELINE_VARIABLE_CUSTOM="MY_CUSTOM_VALUE"
export PIPELINE_VARIABLE_EMAIL="stejon@freepik.com"
export PIPELINE_VARIABLE_REF="main"
export PIPELINE_VARIABLE_SHORTCOMMIT="cdd6128"
export PIPELINE_VARIABLE_TAG="false"
export PIPELINE_VARIABLE_USER="sergiotejon"
export PIPELINE_COMMIT="cdd61288696ace5125d89b03a4de73fefd0f9e3e"
export PIPELINE_DIFF_COMMIT="dac9c537e2a6518bb4a30a455b241d88412bfec1"
export PIPELINE_REPOSITORY="git@github.com:sergiotejon/repo-github.git"
export PIPELINE_EVENT="push"

rm -rf /tmp/repo
go run ../cmd/pipeline-converter/main.go -c ../configs/config_example.yaml

