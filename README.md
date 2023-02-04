# CEC Cascade

This is a little home automation tool that I've written to help manage my HiFi system, which has no "smart" features like HDMI (and therefore eARC), nor sensitivity to incoming signals that are all done via S/PDIF.

In my particular system, there are three possible "modes" for consumption. Central to this HiFi setup is a DAC, which is used to switch between these modes.

## Modes

1. **General music playback (default mode)**:
    * **Source**: AirPort Express
    * **Connection**: Optical S/PDIF
    * **Destination**: DAC Input #4 (optical) 
2. **Critical music listening**:
    * **Source**: RoonBridge on Raspberry Pi (RPi 3)
    * **Connection**: Coax S/PDIF
    * **Destination**: DAC Input #1 (coax)
3. **TV viewing**:
    * **Source**: TV
    * **Connection**: Optical S/PDIF
    * **Destination**: DAC Input #3 (optical)

## Signals to the RPi 3
The TV is also connected to the HDMI port on the RPi 3, which broadcasts the power state of the TV using HDMI-CEC.

This means that the RPi 3 is aware of whether:
1. The TV is on or off (via HDMI-CEC)
2. RoonBridge is using the sound output device (by querying the sound device state)

> **N.B.**: The RPi 3 is _not_ aware if AirPlay is active between the AirPort Express and the DAC. But seeing as modes 2 & 3 are done intensively, if neither are active then the system _should_ default to mode 1.

## Mode priority
The modes above have been listed in ascending order of priority. If we know mode 2 is active, it should override mode 1. If we know mode 3 is active, it should override mode 1 and 2.

|                    	| AirPlay mode 	| RoonBridge mode 	| TV mode 	|
|--------------------	|----------------	|-------------------	|-----------	|
| Source: AirPlay    	| ✅              	| ❌                 	| ❌         	|
| Source: RoonBridge 	| *️⃣              	| ✅                 	| ❌         	|
| Source: TV         	| *️⃣              	| *️⃣                 	| ✅         	|

> ### **Legend**
> ✅ Source is active <br />
> *️⃣ Source status ignored <br />
> ❌ Source must be off

## Dependencies

This app can be built using the `make build` command, which will output a linux/arm64 executable (for use with the RPi 3). There are a few system requirements for this application to work.

### cec-client

The `cec-client` application is a standalone command-line application that is bundled with the `cec-utils` Debian package. For Raspberry Pis, this can be installed using the following command:

```
sudo apt-get install cec-utils
```

More information about this application can be found in the [rmtsrc/raspberry-pi-cec-client.md](https://gist.github.com/rmtsrc/dc35cd1458cd995631a4f041ab11ff74) Gist.

The `cec-cascade` invokes the `cec-client` binary, which _must_ have the privilages to query the HDMI-CEC status (check that the `cec-client` binary is owned by root).

### python-broadlink
The [`python-broadlink` toolset](https://github.com/mjg59/python-broadlink) provides both discovery and command issuing utilities that interface with a Broadlink RM3/4 Pro, which is what this utilty uses to issue infrared signals.

Specifically, the `broadlink_discovery` and `broadlink_cli` python commands must be available in the $PATH.

This application assumes that the device has been set up using the Broadlink app, so that it is already connected to the home WiFi and ready to accept commands.
