# gen-userdata
A utility for generating an Ignition config file from a Golang template.
Makes it trivial to convert permissions to decimal, encode files, and escape systemd unit contents.

## Installation
`go get -u github.com/goeieware/gen-userdata`

## Example of use:
`gen-userdata example/userdata.tmpl > config.ign`

## Functions provided to the template:
* `mode`: Converts supplied file mode/permission, and inserts the decimal equivalent into the ignition config.
* `base64`: Loads the supplied file, encodes as Base64, and inserts the string into the ignition config. 
* `exscape`: Loads the supplied files, escapes quotes and newlines, and inserts the result into the ignition config.

### Example template:

```
{
    "ignition": { 
        "version": "2.0.0"
    },
    "storage": {
        "files": [
            {
                "filesystem": "root",
                "path": "/var/lib/iptables/rules-save",
                "mode": {{ mode 0644 }},
                "owner": "root:root",
                "contents": {
                    "source": "{{ base64 "example/iptables/rules-save" }}"
                }
            },
            {
                "filesystem": "root",
                "path": "/home/core/initialize.sh",
                "mode": {{ mode 0555 }},
                "contents": {
                    "source": "{{ base64 "example/initialize.sh" }}"
                }
            }
        ]
    },
    "systemd": {
        "units": [
            {
                "name": "iptables-restore.service",
                "enable": true
            },
            {
                "name": "initialize.service",
                "user": "core",
                "enable": true,
                "contents": "{{ escape "example/systemd/initialize.service" }}" 
            },
            {
                "name": "mnt-datastore.mount",
                "enable": true,
                "contents": "{{ escape "example/systemd/mnt-datastore.mount" }}" 
            }
        ]
    }
}
```
