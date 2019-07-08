---
layout: post
title: "Re-encrypt LUKS boot partition"
date: 2018-02-23 13:13:58 +0100
categories: luks encryption
redirect_from:
  - /post/reencrypt-luks-boot-partition
---

> There is now a better way to do this using `cryptsetup-reencrypt`.
> See `man cryptsetup-reencrypt` for more information.

I was not happy with the performance I was getting with a key-size of 512, and decided to switch to a length of 256 to reduce time spent decrypting on boot and resume from disk. That balance between security and usability.

Start by booting normally into your OS.

**WARNING** Take a backup of your system before continuing. At the **very least**, make absolutely sure you have all the decryption keys for all devices  backed up on a separate device that you indeed do have un-encrypted access to. I almost lost my decryption key for my LVM device due to having a solution with a decryption file embedded into the initramfs. Take caution!

> It is assumed your system is set up something like in my other post; [LVM on LUKS with encrypted boot partition](https://blog.stigok.com/2018/05/03/lvm-in-luks-with-encrypted-boot-partition-and-suspend-to-disk.html).)

Backup the `/boot` partition (this will include the EFI, but doesn't really matter)

    # tar -cf /root/boot-backup-$(date +%F).tar.gz /boot

Unmount EFI and boot partition

    # umount /boot/efi
    # umount /boot

My partitions are like the following

    nvme0n1   physical disk device
    nvme0n1p1 EFI partition
    nvme0n1p2 Encrypted boot partition
    nvme0n1p3 Encrypted LVM partition (root, swap and home)

Reformat the boot device. (Yes, of course, **make absolutely sure** you are selecting the correct device and partition).
I am re-encrypting my boot partition, so I am going to select `nvme0n1p2`.

    # cryptsetup -v --key-size=256 --hash sha256 --iter-time 5000 --use-random luksFormat /dev/nvme0n1p2

You will be asked for a password. Make sure you enter one which you can remember for certain.

Reopen the boot partition

    # cryptsetup luksOpen /dev/nvme0n1p2 encrypted-boot

In addition to your password, create a keyfile to unlock the boot partition so that the boot partition can be automatically mounted by systemd crypttab on successful boot.

If you want increased entropy in the generation of the keyfile, install `haveged`, make sure it's started, and use `/dev/urandom` instead of `/dev/random` when generating. Creating a 4096 length random key.

    # mkdir -p /etc/initcpio/keyfiles
    # dd bs=512 count=8 iflag=fullblock if=/dev/urandom of=/etc/initcpio/keyfiles/encrypted-boot.key

Add the keyfile as a decryption key of the volume

    # cryptsetup luksAddKey /dev/nvme0n1p2 /etc/initcpio/keys/encrypted-boot.key

Verify that you used the correct key. Make sure the below command doesn't emit a message like `No key available with this passphrase.`

    # cryptsetup open /dev/nvme0n1p2 encrypted-boot-open-test --key-file /etc/initcpio/keys/encrypted-boot.key 
    Cannot use device /dev/nvme0n1p2 which is in use (already mapped or mounted)

If the above looks alright, format then mount the boot volume

    # mkfs.ext2 /dev/mapper/encrypted-boot
    # mount /dev/mapper/encrypted-boot /boot

Re-mount EFI partition on top

    # mount /dev/nvme0n1p1 /boot/efi

Verify that the UUID of the boot partition corresponds with the GRUB config

    # blkid -s UUID /dev/nvme0n1p3
    /dev/nvme0n1p3: UUID="31c1fadf-4a5c-40a5-b1af-05160db54699"

    # grep 'GRUB_CMDLINE_LINUX=' /etc/default/grub
    GRUB_CMDLINE_LINUX="cryptdevice=UUID=31c1fadf-4a5c-40a5-b1af-05160db54699:encrypted-lvm root=/dev/mapper/Main-root resume=/dev/mapper/Main-swap cryptkey=rootfs:/encrypted-lvm.key"

Note the path of the `cryptkey=`. The key file to decrypt the **LVM** (root) partition will be embedded into the initramfs image.

Make sure the `CRYPTODISK` option is enabled

    # grep CRYPTODISK /etc/default/grub 
    GRUB_ENABLE_CRYPTODISK=y

Re-generate GRUB config and re-install EFI image

    # grub-mkconfig -o /boot/grub/grub.cfg
    # grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=grub --recheck

### Embed the keyfile to decrypt the LVM partition into the initramfs image

Create a new install-script at `/etc/initcpio/install/decryption-keys` which will take care of adding the keyfiles to the image whenever `mkinicpio` is executed.
Files may be manually entered into the `FILES=()` section of `/etc/mkinitcpio.conf` as well, but I it seems to force you into mounting it at the same path within the image as in your local system. This is why I created an install script to embed it from a different local path. (TODO: needs rephrasing).

    cat <<'EODUMP' > /etc/initcpio/install/decryption-keys
    #!/bin.bash
    # https://gist.github.com/stigok/7c8d3c872fae5573a870ecd86a4c896c
    
    KEYDIR=/etc/initcpio/keys
    
    function help {
      cat <<EOF
    This hook will embed decryption keys for the encrypted root device into
    initramfs to automatically mount the root partition after a successful
    decryption of the boot partition.
    Expects keyfiles to reside in $KEYDIR with files named after their mount name
    E.g: $KEYDIR/encrypted-boot.key
    EOF
    }
    
    function build {
      # Add all available keys
      for file in $KEYDIR/*; do
        [ -e "$file" ] || continue
        add_file "$file" "$(basename $file)" 0400
      done
    }
    EODUMP

Reference the script in the `HOOKS` section of `/etc/mkinitcpio.conf`

    # grep 'HOOKS=' /etc/mkinitcpio.conf
    HOOKS=(base udev keyboard keymap consolefont autodetect modconf block encrypt lvm2 resume decryption-keys filesystems fsck)

Place the decryption key of the LVM device in the `KEYDIR` folder and make it the same as the volume name, with a `.key` extension. (TODO: This path should maybe be at  /etc/initcpio` somewhere)
If you don't have a keyfile, or if it's embedded in the old initcpio image and needs to be extracted, see other walkthrough to recover. (TODO: Add link)

    # mv ~/lvm-decryption-key /etc/initcpio/keys/encrypted-lvm.key

Set hard permissions on the keys

    # chown root:root /etc/initcpio/keys/*
    # chmod 0000 /etc/initcpio/keys/*
    # chattr +i /etc/initcpio/keys/*

The rationale for putting these files here is that they are residing within the encrypted LVM, and hence needs to be decrypted properly in order to be accessed.

Generate the initramfs. Add the `-v` arg for verbosity to see that the keyfiles are being added explicitly.

    # mkinitcpio -p linux -v

Cross your fingers and reboot - in that order.

Get rid of the backup in a secure manner. Remember that the initramfs contains the decryption key of the LVM.