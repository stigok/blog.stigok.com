---
layout: post
title: "My Arch Linux installation notes"
date: 2017-12-04 11:33:59 +0100
categories: arch draft
redirect_from:
  - /post/my-arch-linux-installation-notes
---

All the links I visited for information about my preferred setup

- https://wiki.archlinux.org/index.php/Installation_guide
- https://tothepoles.wordpress.com/2017/11/16/lenovo-t470p-ubuntu-16-04-install-notes/
- https://wiki.archlinux.org/index.php/Dm-crypt/Encrypting_an_entire_system
  - https://superuser.com/questions/763642/secure-erase-ssd-on-lenovo-thinkpad-t520-cant-unfreeze-ssd-machine-reboots-on
  - https://pcsupport.lenovo.com/no/en/downloads/ds019026
  - https://unix.stackexchange.com/questions/257363/create-a-uefi-bootable-usb-stick-from-an-iso-image-with-gparted

## LVM on LUKS with encrypted boot partition using Grub

Create en EFI System Partition (ESP). I'm setting this as 512MB to avoid potential disk space issues in the years to come.

    # gdisk /dev/nvme0n1
    o (new partition table)
    y (confirm)
    n (new partition)
    [blank] (default partition number)
    [blank] (default start sector)
    +512M (last sector)
    ef00 (BIOS boot)

While at it, create the boot and LVM partition as well without leaving `parted`.

    n (new partition)
    [blank] (default partition number)
    [blank] (default start sector)
    +1G (last sector)
    8300 (Linux filesystem)

    n (new partition)
    [blank] (default partition number)
    [blank] (default start sector)
    [blank] (last sector (all available space))
    8e00 (Linux LVM)

    w (write changes to disk)
    y (confirm write)

Create a FAT32 file system on the EFI system and boot partition

    # mkfs.vfat -F 32 /dev/nvme0n1
    # mkfs.vfat -F 32 /dev/nvme0n2

Verify configuration up until now by comparing the following command and results

    # parted /dev/nvme0n1 print
    [output pending]

Create LUKS container on system partition. You will be asked for confirmation, then prompted for a passphrase when running this command. **Make sure this passphrase is memorable to you**.

    # cryptsetup -v --key-size 512 --hash sha512 --iter-time 5000 --use-random luksFormat

Decrypt, open and mount the encrypted container. Here, *cryptolvm* is the name of which the container will be mapped as.

    # cryptsetup open /dev/nvme0n1p3 cryptolvm

Create physical volume under encrypted LVM

    # pvcreate /dev/mapper/cryptolvm

Create a volume group under this volume. *Main* is an arbitrary name I chose for the group

    # vgcreate Main /dev/mapper/cryptolvm

Create wanted logical volumes on the volume group. The swap partition should be greater than or equal to the amount of system RAM for suspend-to-disk support. [According to ArchWiki](https://wiki.archlinux.org/index.php/Power_management/Suspend_and_hibernate#About_swap_partition.2Ffile_size), hibernation *may be successful* even if swap partition is smaller, but I want to make sure it **always** succeeds.

    # lvcreate -L 16G Main -n swap

My existing system has 25GB root, but <15GB use after a year without care for. However, your milage may vary greatly depending on your usage.

    # lvcreate -L 25G Main -n root

Using 200GB for the home partition instead of `100%FREE` to leave space for other operating systems I might want to install later on.

    # lvcreate -L 200G Main -n home

Create filesystems on each logical volume

    # mkfs.ext4 /dev/mapper/Main-root
    # mkfs.ext4 /dev/mapper/Main-home
    # mkswap /dev/mapper/Main-swap

Prepare the boot partition outside of LVM

    # mkfs.ext2 /dev/nvme0n1p2

## Mount the logical volumes

    # mount /dev/mapper/Main-root /mnt
    # mkdir /mnt/home
    # mount /dev/mapper/Main-home /mnt/home
    # swapon /dev/mapper/Main-swap

Create `boot` folder on the root volume and mount the boot partition

    # mkdir /mnt/boot
    # mount /dev/nvme0n1p2 /mnt/boot

Mount the ESP

    # mkdir /mnt/boot/efi
    # mount /dev/nvme0n1p1 /mnt/boot/efi

## Prepare for chroot

Re-order the list of package mirrors pacman is using if desired. I am ordering them by in `/etc/pacman.d/mirrorlist`. This file will be copied to the new system.

If you're plugged in with ethernet but haven't configured an IP address yet, try get one through DHCP before continuing to the next step.

    # dhclient enp0s31f6

Install base packages and all other things you might consider *essential* to you.

    # pacstrap /mnt base base-devel vim grub efibootmgr

Generate fstab file. Review it afterwards to make sure it looks right.

    # genfstab -U /mnt >> /mnt/etc/fstab

Change root into the new system

    # arch-chroot /mnt

Set the time zone and sync the hardware clock, assuming it is set to UTC.

    # ln -sf /usr/share/zoneinfo/Europe/Oslo /etc/localtime
    # hwclock --systohc

Configure desired locales in `/etc/locale.gen` and generate with `locale-gen`, then configure defaults in `/etc/locale.conf` afterwards. I want English language with Norwegian date and time:

    # cat > /etc/locale.conf
    LANG=en_US.UTF-8
    LC_TIME=nb_NO.UTF-8

## mkinitcpio

Configure `/etc/mkinitcpio.conf` to use encrypt with systemd. Make sure you have configured your keymap in `/etc/vconsole.conf`, or you might have a terrible time entering your decryption passphrase when booting. 

I am using `keyboard` before `autodetect` to make sure all keyboard drivers are loaded. If an external keyboard is connected later on (e.g. by docking) it will not have a driver available s.

    HOOKS=(base systemd keyboard autodetect sd-vconsole modconf block sd-encrypt sd-lvm2 filesystems fsck)

Create the initial ramdisk environment

    # mkinitcpio -p linux

- https://wiki.archlinux.org/index.php/Mkinitcpio#Common_hooks

## Create bootloader with GRUB

Get GRUB installed first of all

    # pacman -S grub

Now update the following line in `/etc/default/grub`

    GRUB_CMDLINE_LINUX="luks.uuid=%uuid%"

Then replace `%uuid%` with the part UUID of the **block device containing** LUKS. This can of course be done manually, but from terminal only it easier to do with `sed`

    # sed -i s/%uuid%/$(blkid -o value -s PTUUID /dev/nvme0n1)/ /etc/default/grub

- https://wiki.archlinux.org/index.php/Dm-crypt/Encrypting_an_entire_system#Preparing_the_logical_volumes

Generate GRUB configuration

    # grub-mkconfig -o /boot/grub/grub.cfg

Verify that the ESP is mounted to `/boot/efi` with `lsblk`.

Install bootloader to the ESP

    # grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=grub --recheck

- https://wiki.archlinux.org/index.php/Dm-crypt/System_configuration#mkinitcpio