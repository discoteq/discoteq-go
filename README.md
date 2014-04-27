Discoteq
========

`discoteq` is a tool for service discovery integration. Its job is to
attach your application to all those fancy service coördination systems
that are popping up. In an ideal world, you wouldn't need discoteq. But
when you've got to iterate into something better, discoteq can help.

The simplist thing that could possibly work
-------------------------------------------

    apt-get install discoteq
    cat >/etc/discoteq-chef.json <<EOF
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
    discoteq-chef &
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

This assumes the existance of a usable `/etc/chef/client.rb` with the necessary credentials available.

Installation
------------

    apt-get install discoteq


Usage
-----

explicit config file

explicit output file

change event triggers




Contributing
-----

go get

https://github.com/marpaia/chef-golang.git