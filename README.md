# debug CNS registration of in-tree volumes

This tool allows us to debug CNS registration of in-tree vSphere volumes without
potentially enabling CSI migration.

## How to use it.

### Build it:

```
~> make
```
### Run it:

```
~> export KUBECONFIG=<path_to_readable_kubeconfig>
~> ./bin/cns-register -pv pvc-b81a0f6b-aeaf-42d6-9832-2c9c912d9c18
```

Where `pvc-b81a0f6b-aeaf-42d6-9832-2c9c912d9c18` is the name of the PV you are looking to register with vCenter.


This tool only works in OCP clusters.
