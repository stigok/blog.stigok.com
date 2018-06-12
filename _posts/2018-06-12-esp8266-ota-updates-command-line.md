---
layout: post
title:  "Flash ESP8266 firware over the air (OTA) with espota.py"
date:   2018-06-13 00:18:05 +0200
categories: esp8266
---

## Introduction

I have been spending a lot of time getting a setup I've been satisfied with
for flashing my ESP-12F, but every time I improve it, I want to improve
further.

Now I want to avoid having to physically connect to the ESP when I want to
flash, and instead send the new image through the wireless network.

This is a continuation of my previous post about [flashing using `esptool`][previous-post].

## Download espota.py

The [*espota.py*][espota.py] is a pretty small script that takes care of OTA updates.
Download and set execute permissions for easy usage

    $ curl -O https://github.com/esp8266/Arduino/blob/master/tools/espota.py
    $ chmod +x espota.py

Now move that file to a folder in your path.

## Prepare an OTA image

This first time, I have to be physically connected. When the OTA code has
been added, it will listen for incoming connections on port 8266 by default,
and expect a payload containing the new image to flash itself with.

Take a look at this snippet taken from the [esp8266/Arduino][repo] GitHub repo
([BasicOTA.ino][BasicOTA.ino])

```cpp
#include <ESP8266WiFi.h>
#include <ESP8266mDNS.h>
#include <WiFiUdp.h>
#include <ArduinoOTA.h>

const char* ssid = "..........";
const char* password = "..........";

void setup() {
  Serial.begin(115200);
  Serial.println("Booting");
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  while (WiFi.waitForConnectResult() != WL_CONNECTED) {
    Serial.println("Connection Failed! Rebooting...");
    delay(5000);
    ESP.restart();
  }

  // Port defaults to 8266
  // ArduinoOTA.setPort(8266);

  // Hostname defaults to esp8266-[ChipID]
  // ArduinoOTA.setHostname("myesp8266");

  // No authentication by default
  // ArduinoOTA.setPassword("admin");

  // Password can be set with it's md5 value as well
  // MD5(admin) = 21232f297a57a5a743894a0e4a801fc3
  // ArduinoOTA.setPasswordHash("21232f297a57a5a743894a0e4a801fc3");

  ArduinoOTA.onStart([]() {
    String type;
    if (ArduinoOTA.getCommand() == U_FLASH)
      type = "sketch";
    else // U_SPIFFS
      type = "filesystem";

    // NOTE: if updating SPIFFS this would be the place to unmount SPIFFS using SPIFFS.end()
    Serial.println("Start updating " + type);
  });
  ArduinoOTA.onEnd([]() {
    Serial.println("\nEnd");
  });
  ArduinoOTA.onProgress([](unsigned int progress, unsigned int total) {
    Serial.printf("Progress: %u%%\r", (progress / (total / 100)));
  });
  ArduinoOTA.onError([](ota_error_t error) {
    Serial.printf("Error[%u]: ", error);
    if (error == OTA_AUTH_ERROR) Serial.println("Auth Failed");
    else if (error == OTA_BEGIN_ERROR) Serial.println("Begin Failed");
    else if (error == OTA_CONNECT_ERROR) Serial.println("Connect Failed");
    else if (error == OTA_RECEIVE_ERROR) Serial.println("Receive Failed");
    else if (error == OTA_END_ERROR) Serial.println("End Failed");
  });
  ArduinoOTA.begin();
  Serial.println("Ready");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  ArduinoOTA.handle();
}
```

This script doesn't do much other than updating itself, and printing the
progress to the serial output. But this is a great first script for
testing that it actually works.

## Build image

Update the `ssid` and `password` string variables before building the file
(for information about *build.options.json*, see my [previous esp 8266 post][previous-post])

```terminal
$ arduino-builder -compile -logger=humantags -warnings=default \
                  -build-cache=$(pwd)/.cache -verbose \
                  -build-options-file=$(pwd)/build.options.json \
                  -build-path=$(pwd)/build \
                  $(pwd)/ota_test.ino
```

This should produce a binary image file *ota_test.ino.bin*

## Upload

Upload the image to the *physically connected* esp8266

    $ esptool.py --baud 115200 -port /dev/serial0 write_flash --flash_mode qio 0x00000 ota_test.ino.bin

I can now use connect to the serial interface to watch its output while I try
to update over the air. The chip should print its IP address to the serial
output when it has connected successfully to the wireless network.
(I use `minicom -b 9600 -D /dev/serial0` right now)

Update itself with the exact same image to verify that it's working as
intended. I use the IP address it has outputted to the serial output upon boot

    $ espota.py -i 10.10.3.71 -f ../build/ota_test.ino.bin

![minicom esp uploading ota](https://public.stigok.com/img/1528841616419416938.png)
![minicom esp successful ota update](https://public.stigok.com/img/1528841728778692471.png)

## Update button trigger

Next up is adding code to read a GPIO pin connected to a button, to only call
`ArduinoOTA.handle()` when a button is pressed.

For example by putting this in `setup()`:

    pinMode(BUTTON_PIN, INPUT);

a `loop` definition similar to the following:

```cpp
void loop() {
  Serial.printf("button state: %u\r", digitalRead(BUTTON_PIN));

  const unsigned short int allowOtaUpdate = digitalRead(BUTTON_PIN);
  if (allowOtaUpdate) {
    ArduinoOTA.handle();
  }

  delay(1000);
}
```

Defining the `BUTTON_PIN` like `const unsigned short int BUTTON_PIN = 5;` and
connect a button to *GPIO05* on the chip.

I am connecting this chip to a thermal receipt printer, so I will be able to
read the serial output even though I'm not pysically connected to it.

## References
- https://www.elecrow.com/download/ESP-12F.pdf

[previous-post]: https://blog.stigok.com/2018/06/01/build-arduino-sketches-from-command-line.html
[repo]: https://github.com/esp8266/Arduino
[BasicOTA.ino]: https://github.com/esp8266/Arduino/blob/03baea27efddb819ba15ffed9b47ee9da8410f54/libraries/ArduinoOTA/examples/BasicOTA/BasicOTA.ino
[espota.py]: https://github.com/esp8266/Arduino/blob/master/tools/espota.py)
