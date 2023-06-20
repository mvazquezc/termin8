# termin8

CLI tool that cleans finalizers from stuck kubernetes objects.

~~~console
Terminates stuck namespaced resources in the specified namespaces

Usage:
  termin8 run [flags]

Flags:
  -d, --dry-run                      Will not terminate stuck resources, will output what would have been terminated
  -o, --extended-output string       Extended output in an specific format. Usage: '-o [  yaml | json ]'
  -h, --help                         help for run
  -k, --kubeconfig string            Path to the kubeconfig file to be used. If not set, will default to KUBECONFIG env var
  -n, --namespaces strings           List of namespaces where stuck objects will be terminated (comma separated) e.g: ns1,ns2
  -s, --skip-api-resources strings   List of namespaced api resources to skip (comma separated) e.g: myresource.group.example.com,myresource2.group2.example.com
~~~

## How does it work?

The tool will consume the kubeconfig file passed as an argument or will try to read the `KUBECONFIG` env var to get one.
Once connected it will get all namespaced API resources available in the cluster and will get the ones being deleted and
with pending finalizers in the namespaces passed as an argument.

Once the complete list of resources is ready, the finalizers will be cleaned by the tool.

## Example runs

### Clean stuck resources on namespace default

~~~shell
$ termin8 run -n default

✓ termin8 completed. 10 stuck resources terminated.
~~~

### Clean stuck resources on namespace default and output extended yaml

~~~shell
$ termin8 run -n default -o yaml

✓ termin8 completed. 10 stuck resources terminated.

- namespace: default
  terminatedresources:
    - crontabs/stuck-object-2jxz2
    - crontabs/stuck-object-8sf75
    - crontabs/stuck-object-b5f9j
    - crontabs/stuck-object-cjszg
    - crontabs/stuck-object-gwlt9
    - crontabs/stuck-object-kg5ms
    - crontabs/stuck-object-l24wh
    - crontabs/stuck-object-l98ps
    - crontabs/stuck-object-mrvcs
    - crontabs/stuck-object-x9tmf
~~~

### Dry-run for stuck resources on namespace default

~~~shell
$ termin8 run -n default --dry-run

✓ termin8 completed. 10 stuck resources would have been terminated.

- namespace: default
  terminatedresources:
    - crontabs/stuck-object-7vxk5
    - crontabs/stuck-object-bz4ml
    - crontabs/stuck-object-gng5s
    - crontabs/stuck-object-j869v
    - crontabs/stuck-object-l82kv
    - crontabs/stuck-object-lgvb2
    - crontabs/stuck-object-rr62m
    - crontabs/stuck-object-vgzrd
    - crontabs/stuck-object-w7vm4
    - crontabs/stuck-object-xlm29
~~~


## Important

This tool is to be used as a last resort, you shouldn't be removing finalizers. Using this tool may lead to a broken cluster,
orphaned cloud resources, etc.

## Future work

- Add tests
