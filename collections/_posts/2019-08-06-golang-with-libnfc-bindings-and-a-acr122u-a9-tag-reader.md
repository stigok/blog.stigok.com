---
layout: post
title:  "Using libnfc bindings in Golang with an ACR122U-A9 tag reader"
date:   2019-08-06 20:00:45 +0200
categories: nfc golang
---

I'm playing with an NFC tag reader (ACR122) and libnfc in Go.
Getting the reader to work properly in Linux is a bit of a task in itself.

## Make reader work with libnfc on Linux

For the current version of the Go bindings (2.1.4) you need libnfc v1.8.0 or later.
libnfc is probably in your package manager's default repositories, but note the
version you are offered.

Since this article focues on Go, more details for the Linux setup
can be read at the [Arch Linux Wiki](https://wiki.archlinux.org/index.php/Touchatag_RFID_Reader).

## Allow unprivileged device access

This will allow access to the device without the need for root user or `sudo`.
Create a new group for the nfc device:

```
# groupadd --system nfc
```

Add your user to the `nfc` group:

```
$ sudo gpasswd --add myusername nfc
```

Determine the vendor and product id of USB device (`ID <vendor>:<product>`).
Look for a device that has something to do with NFC, cards or the product
name of your NFC reader.

```
$ sudo lsusb
Bus 001 Device 007: ID 072f:2200 Advanced Card Systems, Ltd ACR122U
[redacted output]
```

Create a new `udev` (dynamic device management in Linux) rule to allow the `nfc` group read and write access to the device.
This rule will make sure specific permissions are set on all devices that matches the vendor and product id we found in the last step.

```
$ echo SUBSYSTEMS=="usb", ATTRS{idVendor}=="072f", ATTRS{idProduct}=="2200", SYMLINK+="nfc/acr122", GROUP="nfc", MODE="660" | sudo tee /etc/udev/rules.d/99-libnfc.rules
```

Close your terminal window and reopen it to refresh your group ownership.
Connect your device to the computer, or reconnect it if it is already connected.
The device should now be accessible by all users part of the `nfc` group.

## Test reader connection in Go program

I'm using [clausecker's NFC bindings][nfc-repo]. I had prior experience with `libnfc`,
so I wanted to keep things simple for myself. The library reference can be [found here](https://pkg.go.dev/github.com/clausecker/nfc/v2).

I had to search for other projects using this library as a dependency to get a better
grasp of how it should be used. I used [arsatiki/pocket-gopher][] as a reference implamentation.

Install the dependency while inside the Go project folder. If you are using go modules,
this will save a reference to this dependency in your go.mod file.

```
$ go get github.com/clausecker/nfc/v2
```

I can now use it in my program source code. The below code is how I interfaced with my
NFC device.

```go
package main

import (
	"fmt"
	"log"

	nfc "github.com/clausecker/nfc/v2"
)

var (
	// These settings works with the ACR122U. Your milage may vary with
	// other devices.
	m = nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}
	// Use an empty string to select first device libnfc sees
	devstr = ""
)

// This will detect tags or cards swiped over the reader.
// Blocks until a target is detected and returns its UID.
// Only cares about the first target it sees.
func GetCard(pnd *nfc.Device) ([10]byte, error) {
	for {
		targets, err := pnd.InitiatorListPassiveTargets(m)
		if err != nil {
			return [10]byte{}, fmt.Errorf("failed to list nfc targets: %w", err)
		}

		for _, t := range targets {
			if card, ok := t.(*nfc.ISO14443aTarget); ok {
				return card.UID, nil
			}
		}
	}
}

func main() {
	log.Println("using libnfc", nfc.Version())

	pnd, err := nfc.Open(devstr)
	if err != nil {
		log.Fatalln("could not open device:", err)
	}
	defer pnd.Close()

	if err := pnd.InitiatorInit(); err != nil {
		log.Fatalln("could not init initiator:", err)
	}

	log.Println("opened device", pnd, pnd.Connection())

	card_id, err := GetCard(&pnd)
	if err != nil {
		log.Fatalln("failed to get_card", err)
	}

	if card_id != [10]byte{} {
		// Print card ID as uppercased hex
		log.Printf("card found: %#X", card_id)
	} else {
		log.Println("no card found")
	}
}
```

Example program output:

```
$ go build nfc.go
$ ./nfc
2020/10/29 19:27:34 using libnfc 1.8.0
2020/10/29 19:27:34 opened device ACS / ACR122U PICC Interface acr122_usb:001:012
2020/10/29 19:27:39 card found: 0X044A77B1001620000000
```

## References
- <https://github.com/nfc-tools/libnfc/issues/458>
- <https://wiki.archlinux.org/index.php/Touchatag_RFID_Reader>
- <https://github.com/nfc-tools/libnfc>
- <https://github.com/nfc-tools/libnfc/blob/master/README.md#acr122>
- <https://github.com/fuzxxl/nfc>
- <https://github.com/arsatiki/pocket-gopher/blob/master/list/main.go>

[nfc-repo]: https://github.com/clausecker/nfc
[arsatiki/pocket-gopher]: https://github.com/arsatiki/pocket-gopher/blob/master/list/main.go
