# volumectl-go

Basic go CLI app that uses amixer under the hood to manage audio volume on linux. This is a port of my rust package called [volumectl](https://github.com/blanktiger/volumectl).

# Available arguments

- `-g`/`--get` for displaying current volume in %
- `-i`/`--inc` <VALUE> increase volume by VALUE 
- `-d`/`--dec` <VALUE> decrease volume by VALUE 
- `-t`/`--toggle-mute` toggle mute
