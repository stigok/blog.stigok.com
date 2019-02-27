---
layout: post
title: "Build Arduino sketches from command-line and flash to ESP8266 from Pi"
date: 2018-06-01 19:59:36 +0200
categories: esp8266 arduino raspberrypi
---

I want to build sketches from the command line so that I don't have to use
the Arudino IDE to compile. Then I want to use the `esptool` to flash the
files to the ESP manually. My ESP is connected to a remote Raspberry Pi,
on my desk, so this is a good step towards a faster build process where I
can build on my (fast) computer, then transfer the files to the Pi for it
to my ESP 12-F chip.

## Perform initial Arduino installation

Find your board FQBN by searching in the https://github.com/esp8266/Arduino/blob/master/tools/boards.txt.py
files. My NODEMCU 1.0 (ESP-12E Module) FQBN is `esp8266:esp8266:nodemcuv2`

Install Arduino IDE on your main machine

    # pacman -S arduino

Start the Arduino IDE, go into preferences and paste the ESP8266 URL for "Additional Boards URL"
`http://arduino.esp8266.com/stable/package_esp8266com_index.json` and click Ok.

Open the Board Manager and search for *esp8266*. Install that package, now exit Arduino IDE. 

I've installed `arduino-builder` from the official Arch repo

    # pacman -S arduino-builder

Crate a new project directory somewhere and change to it

    $ mkdir -p ~/esp-project/blink; cd ~/esp-project/blink

Create some directories to house the different files

    $ mkdir src build .cache

Create the test file to make the LED blink every `777` millisecond

    $ cat <<EOF > src/led_blink.ino
    const short int LED = 2; // GPIO 2
    const short int DELAY = 777; // ms

    void setup() {
      pinMode(LED, OUTPUT);
    }

    void loop() {
      digitalWrite(LED, LOW);
      delay(DELAY);
      digitalWrite(LED, HIGH);
      delay(DELAY);
    }
    EOF

Create a *build.options.json* file to simplify the build command

    $ cat <<EOF > build.params.json
    {
      "fqbn": "esp8266:esp8266:nodemcuv2:CpuFrequency=80,FlashSize=4M1M,LwIPVariant=v2mss536,Debug=Disabled,DebugLevel=None____,FlashErase=none,UploadSpeed=115200",
      "builtInLibrariesFolders": "",
      "hardwareFolders": "/usr/share/arduino/hardware,${HOME}/.arduino15/packages",
      "toolsFolders": "/usr/share/arduino/tools-builder,${HOME}/.arduino15/packages",
      "otherLibrariesFolder": "${HOME}/Arduino/libraries",
      "customBuildProperties": "build.warn_data_percentage=75,runtime.tools.mkspiffs.path=${HOME}/.arduino15/packages/esp8266/tools/mkspiffs/0.2.0,runtime.tools.esptool.path=${HOME}/.arduino15/packages/esp8266/tools/esptool/0.4.13,runtime.tools.xtensa-lx106-elf-gcc.path=${HOME}/.arduino15/packages/esp8266/tools/xtensa-lx106-elf-gcc/1.20.0-26-gb404fb9-2",
      "additionalFiles": "",
      "runtime.ide.version": "10805"
    }
    EOF

I can now build my project

    $ arduino-builder -compile -warnings=default -verbose -logger=humantags -build-options-file=$(pwd)/build.options.json -build-path=$(pwd)/build -build-cache=$(pwd)/.cache $(pwd)/src/led_blink.ino

## Transfer the files to the Pi

I'm connected over SSH, so I'm going to use `rsync` to transfer the files over.

    $ rsync build/led_blink.ino.bin pi@10.0.0.42:.

## Flash the chip using `esptool`

** This section is performed on the Raspberry Pi**

Note that I'm running the latest raspbian as of June 1st 2018, which means `/dev/serial0` is a symlink to whichever serial is available for me out of `/dev/ttyAMA0` and `/dev/ttyS0`. If you're on an older system, mind the device path in the commands.

Make sure you have disabled login interface being available on `/dev/serial0`. This can be toggled in the interface options of `raspi-config`.
Failing to do so will give you strange errors when using `esptool`, like invalid headers, failing to write to RAM, or connection failures in general.

I'm not using the raspbian `apt` repo for `esptool`, as the version there is outdated. Instead, download `esptool` using `pip` (`pip3` since it's the python3 `pip` we want) for this. For this post, I'm using version 2.3.1 as it was the latest one when I wrote this.

    $ sudo apt update
    $ sudo apt install python3-pip
    $ sudo pip3 install esptool==2.3.1

Find the *led_blink.ino.bin* file that was compiled and transfered in the previous step and specify that in the flash command

    $ esptool.py --baud 115200 --port /dev/serial0 write_flash --flash_mode qio 0x00000 led_blink.ino.bin
    esptool.py v2.4.0-dev
    Connecting........_____.....____
    Detecting chip type... ESP8266
    Chip is ESP8266EX
    Features: WiFi
    Uploading stub...
    Running stub...
    Stub running...
    Configuring flash size...
    Auto-detected Flash size: 4MB
    Flash params set to 0x0040
    Compressed 250400 bytes to 182663...
    Wrote 250400 bytes (182663 compressed) at 0x00000000 in 16.4 seconds (effective 122.1 kbit/s)...
    Hash of data verified.

    Leaving...
    Hard resetting via RTS pin...

When I've been using the Arduino IDE manually, the chip has reset by itself. Now, however, I have to reset in manually either by pulling the plug, or pushing the button on my prototype board.

## References
- https://github.com/arduino/arduino-builder/blob/master/test/store_build_options_map_test.go
- https://gist.github.com/shadowandy/af468d38ed7ab3ff718b
- https://blog.the-jedi.co.uk/tag/esp8266/

