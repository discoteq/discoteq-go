% DISCOTEQ(1) discoteq | User Commands
% Joseph Holsten
% 2014-10-22

# NAME

discoteq -- discover services from service registries

# SYNOPSIS

`discoteq` [`-c` *config-file*] [`-k` *chef-key*] [`-s` *chef-server*] \
  [`-E` *chef-environment*] [`-u` *chef-client-name*] [*file*]

# DESCRIPTION

`discoteq(1)` is a tool for service discovery integration. Its job is to
attach your application to all those fancy service coordination systems
that are popping up. In an ideal world, you wouldn't need discoteq. But
when you've got to iterate into something better, discoteq can help.

Unlike most other tools in the space, `discoteq(1)` does not handle
scheduling, templating, triggering reloads, or other complex behaviour.
Other tools exist for these problems, notably `cron(8)`, `tilt(1)`,
`kill(1)`, &c.

# OPTIONS

`-c` *config-file*
:   config file path, default `"/etc/discoteq.json"`

`-E` *chef-environment*
:   Chef environment query scope, default:
    `node.chef_environment || "_default"`

`-k` *chef-key*
:   Chef client key file path, default: `"/etc/chef/client.pem"`

`-s` *chef-server*
:   Chef server URL, default: `"http://localhost:4545"`

`-u` *chef-client-name*
:   Chef client username, default: `node.fqdn`


## Supported queries types

Each service query needs to return a list of service host records, each of which must have a `hostname` attribute.

```js
[
  {
    "hostname": "myface-002.example.net",
    "port": 8080
  },
  {
    "hostname": "myface-003.example.net",
    "port": 8081
  }
]
```
As such, queries need a way to select a set of hosts and a way to extract key-value pairs for each.

At the moment, the only type of query supported are Chef searches of the `node` index.

`query`
:   Chef search query.

`role`
:   Shorthand for a `query` of `role:{role}`.

`tag`
:   Shorthand for a `query` of `tag:{tag}`.

`include_chef_environment`
:   Whether to append `"AND chef_environment:{chef-environment}"` to the query. Defaults to true.
    The `chef-environment` comes from the `-E` parameter, or defaults to `"_default"`

`attrs`
:   Map of exported key onto an attribute key. Nested attributes may be accessed by combining the keys with a `.`, eg: setting the `hostname` attribute from `node['cloud']['private_ipv4']` can use the notation `"attrs": { "hostname": "cloud.private_ipv4" }`. `attrs` defaults to `{"hostname": "fqdn"}` and will merge the default with any provided attrs.
 
# EXAMPLES

## The simplest thing that could possibly work

    cat >/etc/discoteq.json <<EOF
    {
        "services": {
            "myface-fascade": {
                "role": "myface-lb"
            },
            "myface": {
                "role": "myface"
            },
            "myface-db-master": {
                "query": "role:myface-db AND tag:master"
            },
            "myface-db-slave": {
                "query": "role:myface-db AND tag:slave"
            },
            "myface-cache": {
                "role": "myface-cache"
            },
            "statsd": {
                "role": "statsd",
                "include_chef_environment": false,
                "attrs": {
                    "hostname": "cloud.private_ipv4",
                    "port": "statsd.port"
                }
            }
        }
    }
    EOF

    discoteq < /etc/discoteq.json > /var/lib/discoteq/services.json
    cat /var/lib/discoteq/services.json
    {
        "services": {
            "myface-fascade": [
                {
                    "hostname": "myface-lb.example.net"
                }
            ],
            "myface": [
                {
                    "hostname": "myface-001.example.net"
                }
            ],
            "myface-db-master": [
                {
                    "hostname": "myface-db-001.example.net"
                }
            ],
            "myface-db-slave": [
                {
                    "hostname": "myface-db-002.example.net"
                },
                {
                    "hostname": "myface-db-003.example.net"
                }
            ],
            "myface-cache": [
                {
                    "hostname": "myface-cache-001.example.net"
                }
            ],
            "statsd": [
                {
                    "hostname": "10.0.0.216",
                    "port": 8126
                }
            ]
        }
    }

## Making it work like `confd`

Run this script through a process manager like `runit` or
`systemd`:

``` {.sh}
#!/bin/sh
# myface-cfg-update - watch for service updates
#   and safely load them into myface

BASE=myface
TEMP=`mktemp -t $BASE.XXXXXXXXXX` || exit 1

SERVICE_MAP=/var/lib/discoteq/services.json
TEMPLATE=/opt/myface/config.erb
CONFIG=/etc/myface.cfg
DAEMON=myface

# watch for service info changes in registry
# re-eval config template
# verify generated config is valid
# swap if valid and different
# reload if swapped
fswatch -1 $SERVICE_MAP &&
  tilt -d $SERVICE_MAP $TEMPLATE > $TEMP &&
  myface --check-cfg $TEMP &&
  flock $CONFIG diffswp $TEMP $CONFIG &&
  service $DAEMON reload
```

With this in your `crontab(5)`:

    */5 * * * *  /usr/local/bin/discoteq \
                 < /etc/discoteq/services-conf.json \
                 > /var/lib/discoteq/services.json

# SEE ALSO

`crontab(5)`, `fswatch(5)`, `tilt(1)`, `flock(1)`,
