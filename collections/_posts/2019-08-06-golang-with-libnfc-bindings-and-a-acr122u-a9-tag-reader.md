---
layout: post
title:  "Using libnfc bindings in Golang with an ACR122U-A9 tag reader"
date:   2019-08-06 20:00:45 +0200
categories: nfc golang
---

I'm playing with an NFC tag reader (ACR122) and libnfc in Go.
Getting the reader to work properly in Linux is a bit of a task in itself.

## Make reader work with libnfc on Linux

Please read the [archlinux wiki page](https://wiki.archlinux.org/index.php/Touchatag_RFID_Reader) for details regarding this step.

## Allow unprivileged (no sudo) device access

Create a new group for the nfc device

```
# groupadd --system nfc
```

Add the user to the group

```
# gpasswd -a myusername nfc
```

Determine vendor and product id of USB device (`ID <vendor>:<product>`)

```
$ lsusb
Bus 001 Device 007: ID 072f:2200 Advanced Card Systems, Ltd ACR122U
```

Create a udev rule to allow `nfc` group read and write access to the device

```
# echo SUBSYSTEMS=="usb", ATTRS{idVendor}=="072f", ATTRS{idProduct}=="2200", SYMLINK+="nfc/acr122", GROUP="nfc", MODE="660" > /etc/udev/rules.d/99-libnfc.rules
```

## Test reader connection in Go program

I'm using [fuzxxl' NFC bindings][nfc-repo]. I picked it because it the most
stars out of the repos I could find, and additionally, since I've used libnfc
before, I wanted to leverage it again. I could not, however, find
any documentation for this library, so I searched GitHub for existing projects
using it as a dependency (see [arsatiki/pocket-gopher][] for more information),
and figured it out from there.

I downloaded the nfc bindings into `$GOPATH/pkg/github.com/fuzxxl/nfc` and
referenced it in my program.

```go
package main

import (
	"fmt"
	"log"

	"github.com/fuzxxl/nfc/2.0/nfc"
)

var (
	m = nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}
	devstr = "acr122_usb:001:012" // Use empty string to select first device
)

// Blocks until a target is detected and returns its UID.
// Only cares about the first target it sees.
func get_card (pnd *nfc.Device) ([10]byte, error) {
	for {
		targets, err := pnd.InitiatorListPassiveTargets(m)
		if err != nil {
			return [10]byte{}, fmt.Errorf("listing available nfc targets", err)
		}

		for _, t := range targets {
			if card, ok := t.(*nfc.ISO14443aTarget); ok {
				return card.UID, nil
			}
		}
	}
}

func main() {
	fmt.Println("using libnfc", nfc.Version())

	pnd, err := nfc.Open(devstr)
	if err != nil {
		log.Fatalf("could not open device: %v", err)
	}
	defer pnd.Close()

	if err := pnd.InitiatorInit(); err != nil {
		log.Fatalf("could not init initiator: %v", err)
	}

	fmt.Println("opened device", pnd, pnd.Connection())

	card_id, err := get_card(&pnd)
	if err != nil {
		fmt.Errorf("failed to get_card", err)
	}

	if card_id != [10]byte{} {
		fmt.Printf("card found %#X\n", card_id)
	} else {
		fmt.Printf("no card found\n")
	}
}
```

Example output of a physical test run

```
$ go build nfc.go
$ ./nfc
using libnfc 1.7.1
opened device ACS / ACR122U PICC Interface acr122_usb:001:012
card found 0X044A77B1001620000000
```

## References
- <https://github.com/nfc-tools/libnfc/issues/458>
- <https://wiki.archlinux.org/index.php/Touchatag_RFID_Reader>
- <https://github.com/nfc-tools/libnfc>
- <https://github.com/nfc-tools/libnfc/blob/master/README.md#acr122>
- <https://github.com/fuzxxl/nfc>
- <https://github.com/arsatiki/pocket-gopher/blob/master/list/main.go>

[nfc-repo]: https://github.com/fuzxxl/nfc
[arsatiki/pocket-gopher]: https://github.com/arsatiki/pocket-gopher/blob/master/list/main.go
