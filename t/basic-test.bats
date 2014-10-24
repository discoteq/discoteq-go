#!/usr/bin/env bats

discoteq=./discoteq

@test "do stuff" {
  nohup goiardi &
  pushd stubs
  knife upload .
  popd
  $discoteq -E dev -c stubs/example.conf |
  jq -e --argfile expected stubs/output.json '. == $expected'
}
