---
layout: post
title: "arch linux gcc 32-bit libraries c assembly"
date: 2017-11-07 23:58:03 +0100
categories: c assembly pacman archlinux
redirect_from:
  - /post/arch-linux-gcc-32bit-libraries-c-assembly
---

I wanted to compile my C code with `gcc -m32` so I could easily follow old C and assembly tutorials, but it required 32-bit gcc libraries which are not available by default in an Arch Linux installation.

From the [ArchWiki](https://wiki.archlinux.org/index.php/Multilib):
> To use the multilib repository, uncomment the [multilib] section in /etc/pacman.conf (Please be sure to uncomment both lines):

After editing my */etc/pacman.conf*, I ran `pacman -Syu` to update the (not yet cached) package list from the multilib repository. Then I installed the package I needed.

    # pacman -S lib32-gcc-libs

Now I can compile my C code for 32-bit machines

    $ gcc -fverbose-asm -g -m32 -O0 -o $dest -c $src

To get the assembly code for the produced object file with Intel notation:

    $ objdump -dS -M intel $dest
    Disassembly of section .text:
    
    00000000 <main>:
    #include <stdio.h>
    
    int main(void)
    {
       0:	8d 4c 24 04          	lea    ecx,[esp+0x4]
       4:	83 e4 f0             	and    esp,0xfffffff0
       7:	ff 71 fc             	push   DWORD PTR [ecx-0x4]
       a:	55                   	push   ebp
       b:	89 e5                	mov    ebp,esp
       d:	53                   	push   ebx
       e:	51                   	push   ecx
       f:	83 ec 10             	sub    esp,0x10
      12:	e8 fc ff ff ff       	call   13 <main+0x13>
      17:	05 01 00 00 00       	add    eax,0x1
      int i = 1;
      1c:	c7 45 f4 01 00 00 00 	mov    DWORD PTR [ebp-0xc],0x1
      printf("%i\n", i);
      23:	83 ec 08             	sub    esp,0x8
      26:	ff 75 f4             	push   DWORD PTR [ebp-0xc]
      29:	8d 90 00 00 00 00    	lea    edx,[eax+0x0]
      2f:	52                   	push   edx
      30:	89 c3                	mov    ebx,eax
      32:	e8 fc ff ff ff       	call   33 <main+0x33>
      37:	83 c4 10             	add    esp,0x10
      return 0;
      3a:	b8 00 00 00 00       	mov    eax,0x0
    }
      3f:	8d 65 f8             	lea    esp,[ebp-0x8]
      42:	59                   	pop    ecx
      43:	5b                   	pop    ebx
      44:	5d                   	pop    ebp
      45:	8d 61 fc             	lea    esp,[ecx-0x4]
      48:	c3                   	ret

## References
- https://stackoverflow.com/q/1289881/90674
- https://stackoverflow.com/q/8021874/90674