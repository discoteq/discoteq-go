#!/usr/bin/env bats

discoteq=./discoteq


setup() {
  goiardi &
}
teardown () {
  kill $(jobs -p)
}

@test "chef services" {
  pushd stubs/chef
  knife upload .
  popd
  $discoteq -E dev -c stubs/chef/example.conf > actual.json
  cat actual.json |
  jq -e --argfile expected stubs/chef/output.json '. == $expected'
  # json-diff actual.json stubs/chef/output.json
}
