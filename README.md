# Discoteq

discoteq -- discover services from service registries

## Synopsis

`discoteq` [`-c` *config-file*] [`-k` *chef-key*] [`-s` *chef-server*]
 [`-E` *chef-environment*] [`-u` *chef-client-name*] [*file*]

## Description

`discoteq(1)` is a tool for service discovery integration. Its job is to
attach your application to all those fancy service coordination systems
that are popping up. In an ideal world, you wouldn't need discoteq. But
when you've got to iterate into something better, discoteq can help.

Unlike most other tools in the space, `discoteq(1)` does not handle
scheduling, templating, triggering reloads, or other complex behaviour.
Other tools exist for these problems: `cron(8)`, `tilt(1)`, `kill(1)`,
&c.

## The simplest thing that could possibly work

    apt-get install discoteq
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

As you can see, you just need to drop a config file, run `discoteq` with
an output file, and it will populate the file with the host attributes
you need.

This assumes the some (somewhat unlikely) defaults for speaking to a
chef server.

## Installation

Eventually, the installation process should be:

    add-apt-repository ppa:discoteq/discoteq
    apt-get update
    apt-get install discoteq

But today you'll need a [working go development environment][], and you
can install with:

    go get github.com/discoteq/discoteq-go

## Usage

By default, `discoteq(1)` expects config as its standard input, and
writes its service map to standard output. Of course, a few things
occasionally need to be specified for your particular environment.

## Options

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
: Chef search query.

`role`
: Shorthand for a `query` of `role:{role}`.

`tag`
: Shorthand for a `query` of `tag:{tag}`.

`include_chef_environment`
: whether to append `"AND chef_environment:{chef-environment}"` to the query. Defaults to true.
  The `chef-environment` comes from the `-E` parameter, or defaults to `"_default"`

`attrs`
: Map of exported key onto an attribute key. Nested attributes may be accessed by combining the keys with a `.`, eg: setting the `hostname` attribute from `node['cloud']['private_ipv4']` can use the notation `"attrs": { "hostname": "cloud.private_ipv4" }`. `attrs` defaults to `{"hostname": "fqdn"}` and will merge the default with any provided attrs.

## Templating from JSON

Of course, your app probably doesn't expect its services in the format
discoteq provides. Perhaps you're currently using chef to populate a
config file using ERB, and it seems awfully inconvenient to have to
change your app to accommodate this tool.

Fear not! The ever powerful [`tilt(1)`][] templating tool will handle
your needs! Simply specify your data and template files and it will
print the rendered text to standard out.

``` {.sh}
discoteq < /etc/discoteq.conf > /etc/myface.json
tilt -d /etc/myface.json myface-template.erb > $APP/myface.config
```

## How do I make it act like [`confd`][]?

[`confd`][] is an awesome tool that tries to solve similar problems to
discoteq. But while `confd` is a one-stop-shop for configuration
management, discoteq requires you to use existing tools to do similar
things.

### How to I schedule discoteq to update files?

At the moment discoteq only supports polling data stores, so scheduling
updates is just a matter of frequency. Ye olde [`crontab(5)`][] is more
than equipped for your needs.

    */5 * * * *  /usr/local/bin/discoteq \
                 < /etc/discoteq/services-conf.json \
                 > /var/lib/discoteq/services.json


### My app doesn't understand discoteq's JSON, how do I populate a template with it?

Your favorite templating language probably already has a command line
client that you can use:

-   [`tilt(1)`][] is a ruby gem that supports a [large number of
    template formats][], and has recently been updated to accept JSON
    input. This is currently our recommended tool for [ERB][]
    templating.
-   [`mustache(1)`][] is the ubiquitous logic-less templating language
    available on every platform ever. This accepts a YAML data input, so
    JSON is fine.

We'd also like to create command-line tools for [jinja][] and
[`test/template`][].

### What if an invalid config file is generated? How do I avoid breaking everything?

Config file check-and-swap is a common problem and deserves a common
solution. At the moment we don't know a good tool to recommend. If you
know of one, please let us know! If you want to build one yourself, we'd
love to help!

### What if someone else is touching the file? Couldn't this lead to inconsistant behaviour?

`flock(1)` allows you to create exclusive locks for writes and shared
locks for reads. It's designed to be used in shell scripts and is
probably easier to use than modifying your existing programs to have
mutual exclusion.

On Linux, you can use the [util-linux `flock(1)`][]. If you aren't on
Linux, we've developed a [portable `flock(1)`][] which supports Darwin,
FreeBSD and Illumos.

### How do I trigger my service when the config updates?

[`fswatch(7)`][] can trigger events whenever a file changes, either as a
one-off or for every event.

    $ fswatch -1 /var/lib/discoteq/services.json |
      xargs -n1 -I% echo "Hello, %!" &
    $ touch /var/lib/discoteq/services.json
    Hello, /var/lib/discoteq/services.json!

Of course this depends triggers on every modification of inode state,
not file contents. Make sure not to touch files you aren't actually
changing!

### That sounds lovely, but what does a real solution look like?

Yeah, tying all this together can be… exciting.

Run this script through a process manager like [`runit`][] or
[`systemd`][]:

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

With this in your [`crontab(5)`][]:

    */5 * * * *  /usr/local/bin/discoteq \
                 < /etc/discoteq/services-conf.json \
                 > /var/lib/discoteq/services.json


## Contributing

Got an idea? Something smell wrong? Cause you pain? Or lost seconds of
your life you'll never get back?

All contributions are welcome: ideas, patches, documentation, bug
reports, complaints, and even something you drew up on a napkin.

Programming is not a required skill. Whatever you've seen about open
source and maintainers or community members saying "send patches or
die" - you will not see that here.

It is more important to me that you are able to contribute.

I promise to help guide this project with these principles:

-   Community: If a newbie has a bad time, it's a bug.
-   Software: Make it work, then make it right, then make it fast.
-   Technology: If it doesn't do a thing today, we can make it do it
    tomorrow.

(Some of the above was repurposed with \<3 from logstash)

For those of you who do want to contribute with code, we've tried to
make it easy to get started. You can install all dependencies and tools
with:

    ./dev-bootstrap.sh

Then you can start support services with:

    forego start

Plenty of example data already exists in `stubs/`, though it probably
deserves more explanation.

Good luck!

  [working go development environment]: http://blog.golang.org/organizing-go-code
  [`tilt(1)`]: https://github.com/rtomayko/tilt/blob/master/man/tilt.1.ronn
  [`confd`]: http://www.confd.io/
  [`crontab(5)`]: http://crontab.org/
  [large number of template formats]: https://github.com/rtomayko/tilt/blob/master/docs/TEMPLATES.md
  [ERB]: https://github.com/rtomayko/tilt/blob/master/docs/TEMPLATES.md#erb
  [`mustache(1)`]: https://mustache.github.io/mustache.1.html
  [jinja]: http://jinja.pocoo.org/
  [`test/template`]: http://golang.org/pkg/text/template/
  [util-linux `flock(1)`]: http://linuxmanpages.net/manpages/fedora20/man1/flock.1.html
  [portable `flock(1)`]: https://github.com/discoteq/flock
  [`fswatch(7)`]: https://github.com/alandipert/fswatch
  [`runit`]: http://smarden.org/runit/index.html
  [`systemd`]: http://freedesktop.org/wiki/Software/systemd/
