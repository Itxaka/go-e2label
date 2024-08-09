# e2label in golang
## Labeling tool that substitutes e2label in pure golang

All done thanks to the upstream gexto lib which created the needed set pieces for this: github.com/nerd2/gexto

This is a simplified and slimmed down version of it and it justs has one purpouse: list and change ext4 labels.

usage is the same as upstream e2label:

```bash
Usage: e2label <filename> [New label]
```

Calling it with a filename or device will return the existing label.
Calling it with a filename or device and a label will set that label.


Thats it, it doesnt need to be more complex than that.

Done thanks to the [Kairos](https://github.com/kairos-io) hackweek of 2024.