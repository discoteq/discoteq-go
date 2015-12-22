#!/usr/bin/env bats

discoteq=./discoteq


setup() {
    consul agent -server \
        -bootstrap-expect=1 \
        -data-dir=/tmp/consul \
        -config-dir=stubs/consul/consul.d &
    sleep 2
}
teardown () {
    kill $(jobs -p)
    rm actual.json
}

@test "consul services" {
  $discoteq -c stubs/consul/discoteq.conf > actual.json
  json-diff stubs/consul/expected.json actual.json
  cat actual.json |
  jq -e --argfile expected stubs/consul/expected.json '. == $expected'
}
