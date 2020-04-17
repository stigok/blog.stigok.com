---
layout: post
title:  "Overriding an existing Perl package in NixOS"
date:   2020-04-16 01:18:35 +0200
categories: nixos nix perl irc irssi
---

I am running irssi as my IRC client. Before I moved from CentOS to NixOS,
I had a hightlighter script that notified me whenever my nickname was mentioned.
This script is written in Perl, just like [other irssi scripts](https://scripts.irssi.org/).

![highlighter, notification-hub and libnotify](https://public.stigok.com/img/2020-04-17-002444.png)

When I migrated to NixOS, the scripts were having some problems.
[scriptassist.pl][] and [highlightcmd][highlighter] were unable to load due to
missing Perl library dependencies.

```
01:23 -!- Irssi v1.2.2 - https://irssi.org
01:23 -!- Irssi: Error in script scriptassist:
01:23 Can't locate LWP/UserAgent.pm in @INC (you may need to install the LWP::UserAgent module) (@INC contains: /home/sshow/.irssi/scripts
          /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/share/irssi/scripts /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/lib/perl5/x86_64-linux-thread-multi
          /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/site_perl/5.30.0/x86_64-linux-thread-multi
          /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/site_perl/5.30.0 /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/5.30.0/x86_64-linux-thread-multi
          /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/5.30.0) at /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/share/irssi/scripts/scriptassist.pl line 24.
01:23 BEGIN failed--compilation aborted at /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/share/irssi/scripts/scriptassist.pl line 24.
```

So then I got started trying to find the names of the packages for the
dependencies I need for this script. Let's start with the first one it nags
about: `perlPackages.LWPUserAgent`.

I try to install it for my current user

```
$ nix-env -iA nixos.perlPackages.LWPUserAgent
```

... then I try to run the script again. But I get the exact same result as
last time. Missing dependency `LWP::UserAgent`.

I try to install it using `sudo`, which will make the package available for
all of the system's users. Maybe that'll make a difference?

```
$ sudo nix-env -iA nixos.perlPackages.LWPUserAgent
```

No. Same result again. Let's read the error message. It's actually
*very* descriptive:

```
Can't locate LWP/UserAgent.pm in @INC (you may need to install the LWP::UserAgent module) (@INC contains: /home/sshow/.irssi/scripts
  /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/share/irssi/scripts /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/lib/perl5/x86_64-linux-thread-multi
  /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/site_perl/5.30.0/x86_64-linux-thread-multi
  /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/site_perl/5.30.0 /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/5.30.0/x86_64-linux-thread-multi
  /nix/store/q9g7sjgicn2hl64c21c8w4h5c0sxjb0w-perl-5.30.0/lib/perl5/5.30.0) at /nix/store/6flnp8hljyg0s520bgj00mmhrpnrbpjg-irssi-1.2.2/share/irssi/scripts/scriptassist.pl line 24.
```

So, `@INC` refers to [paths Perl will look for source files](https://perlmaven.com/how-to-change-inc-to-find-perl-modules-in-non-standard-locations)
when libraries
or modules are imported in a script. When I go through all the
paths listed in the error above, sure enough, there is no `LWP/UserAgent.pm` file to be found.

I run `perl -V` to get a list of all import paths Perl knows about.

```
$ perl -V
  Built under linux
  Compiled at Nov 10 2019 12:37:36
  @INC:
    /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
    /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/site_perl/5.30.1
    /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/5.30.1/x86_64-linux-thread-multi
    /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/5.30.1
```

I obviously need to get some more paths onto that list. Using the Nix REPL I
can determine what library path I want to see in there. I start `nix-repl`
with a channel (string enclosed in `<>`) as an argument, so that I get the `pkgs`
variable available:

```
$ nix repl '<nixpkgs>'
Welcome to Nix version 2.3.3. Type :? for help.

Loading '<nixpkgs>'...
Added 10861 variables.

nix-repl> pkgs.perlPackages.LWPUserAgent.outPath
"/nix/store/rvffh0zclxwfwsv99107dl3zmdqma4cq-perl5.30.0-libwww-perl-6.39
```

I **thought** I wanted to change Perl's environment -- or even recompile Perl
in its entirety to add more default include paths. - Luckily, I found someone
who had already [made a derivation](https://github.com/NixOS/nixpkgs/issues/10350)
for a similar reason. I now came to know some new NixOS swiss army knives:

- [`nixpkgs.config.packageOverrides`](https://nixos.org/nixpkgs/manual/#sec-modify-via-packageOverrides): override derivations for packages
- [`pkgs.<package>.overrideAttrs`](https://nixos.org/nixpkgs/manual/#sec-pkg-overrideAttrs): override derivation attributes for a named package
- [`pkgs.makeWrapper`](https://nixos.org/nixpkgs/manual/#ssec-stdenv-functions): wrap a program to change its runtime environment

All of which came to use in an override derivation for irssi, *irssi-override.nix*:

```nix
# This is an override to support execution of irssi scripts 'scriptassist'
# and 'highlightcmd'
{ nixpkgs, ... }:

{
  # Refer to all packages as nixpkgs, but as pkgs inside the function
  # to avoid infinite recursion.
  nixpkgs.config.packageOverrides = pkgs:
    # Declare common variables for the next derivation
    let
      inherit (pkgs.perlPackages) makePerlPath;
      # Collect all wanted Perl packages in a list, since I'm using it twice
      deps = with pkgs.perlPackages; [
        HTTPDate
        HTTPMessage
        LWP
        LWPUserAgent
        URI
        StringShellQuote
        TextSprintfNamed
        TryTiny
      ];
    in {
      # Override specific attributes in the irssi package
      irssi = pkgs.irssi.overrideAttrs (oldAttrs: {
        # Add all packages as build inputs, including makeWrapper which we
        # will use in the postFixup hook.
        buildInputs = oldAttrs.buildInputs ++ [ pkgs.makeWrapper ] ++ deps;

        # Prepend all of the Perl package's paths to the Perl include path (@INC)
        # using ':' as a string separator
        postFixup = ''
          wrapProgram "$out/bin/irssi" --prefix PERL5LIB : "${makePerlPath deps}"
        '';
      });
    };
}
```

Now, I can include this file in my main file */etc/nixos/configuration.nix* and
add `irssi` to `systemPackages`. (The snippet below has been redacted.)

```
{ config, pkgs, ... }:

{
  imports = [
    ./hardware-configuration.nix
    ./system.nix
    ./network.nix
    ./irssi-override.nix
  ];

  environment.systemPackages = [ pkgs.irssi ];
}
```

I rebuild my system with that configuration and make sure that no parse or
build errors occurs.

```
$ sudo nixos-rebuild test
```

Now, if I'm running `perl -V` **from my local shell**, I will still get the same `@INC`
paths as before, but **`irssi` will have additional paths included in its
environment**. We can verify this from within irssi:

```
/exec perl -V
17:14 Summary of my perl5 (revision 5 version 30 subversion 1) configuration:
[...]
17:14   Compiled at Nov 10 2019 12:37:36
17:14   %ENV:
17:14     PERL5LIB="/nix/store/3x9q1blw2yyavly3li28va3xz5agn0kc-perl5.30.1-HTTP-Date-6.05/lib/perl5/site_perl:/nix/store/f524cgaidn05754i4m1jqvfk7phc8wr3-perl5.30.1-HTTP-Message-6.18/lib/perl5/site_perl:/nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl:/nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl:/nix/store/km34al6vwaw545alf6pvij5drvfk6gbf-perl5.30.1-URI-1.76/lib/perl5/site_perl:/nix/store/8v4accas2f0bn2zk9p5rx0b7n1j9m7dx-perl5.30.1-String-ShellQuote-1.04/lib/perl5/site_perl:/nix/store/78gv0dlk7mhy7jmixw0dd3vwp2f7jxxn-perl5.30.1-Text-Sprintf-Named-0.0403/lib/perl5/site_perl:/nix/store/yv98jyxlhchzdf42cswzgbx9ai8zli68-perl5.30.1-Try-Tiny-0.30/lib/perl5/site_perl"
17:14   @INC:
17:14     /nix/store/3x9q1blw2yyavly3li28va3xz5agn0kc-perl5.30.1-HTTP-Date-6.05/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/3x9q1blw2yyavly3li28va3xz5agn0kc-perl5.30.1-HTTP-Date-6.05/lib/perl5/site_perl/5.30.1
17:14     /nix/store/3x9q1blw2yyavly3li28va3xz5agn0kc-perl5.30.1-HTTP-Date-6.05/lib/perl5/site_perl
17:14     /nix/store/f524cgaidn05754i4m1jqvfk7phc8wr3-perl5.30.1-HTTP-Message-6.18/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/f524cgaidn05754i4m1jqvfk7phc8wr3-perl5.30.1-HTTP-Message-6.18/lib/perl5/site_perl/5.30.1
17:14     /nix/store/f524cgaidn05754i4m1jqvfk7phc8wr3-perl5.30.1-HTTP-Message-6.18/lib/perl5/site_perl
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl/5.30.1
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl/5.30.1
17:14     /nix/store/9j6rvs0c6pkzwqqwj11gvhvyvm66sdmn-perl5.30.1-libwww-perl-6.43/lib/perl5/site_perl
17:14     /nix/store/km34al6vwaw545alf6pvij5drvfk6gbf-perl5.30.1-URI-1.76/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/km34al6vwaw545alf6pvij5drvfk6gbf-perl5.30.1-URI-1.76/lib/perl5/site_perl/5.30.1
17:14     /nix/store/km34al6vwaw545alf6pvij5drvfk6gbf-perl5.30.1-URI-1.76/lib/perl5/site_perl
17:14     /nix/store/8v4accas2f0bn2zk9p5rx0b7n1j9m7dx-perl5.30.1-String-ShellQuote-1.04/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/8v4accas2f0bn2zk9p5rx0b7n1j9m7dx-perl5.30.1-String-ShellQuote-1.04/lib/perl5/site_perl/5.30.1
17:14     /nix/store/8v4accas2f0bn2zk9p5rx0b7n1j9m7dx-perl5.30.1-String-ShellQuote-1.04/lib/perl5/site_perl
17:14     /nix/store/78gv0dlk7mhy7jmixw0dd3vwp2f7jxxn-perl5.30.1-Text-Sprintf-Named-0.0403/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/78gv0dlk7mhy7jmixw0dd3vwp2f7jxxn-perl5.30.1-Text-Sprintf-Named-0.0403/lib/perl5/site_perl/5.30.1
17:14     /nix/store/78gv0dlk7mhy7jmixw0dd3vwp2f7jxxn-perl5.30.1-Text-Sprintf-Named-0.0403/lib/perl5/site_perl
17:14     /nix/store/yv98jyxlhchzdf42cswzgbx9ai8zli68-perl5.30.1-Try-Tiny-0.30/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/yv98jyxlhchzdf42cswzgbx9ai8zli68-perl5.30.1-Try-Tiny-0.30/lib/perl5/site_perl/5.30.1
17:14     /nix/store/yv98jyxlhchzdf42cswzgbx9ai8zli68-perl5.30.1-Try-Tiny-0.30/lib/perl5/site_perl
17:14     /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/site_perl/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/site_perl/5.30.1
17:14     /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/5.30.1/x86_64-linux-thread-multi
17:14     /nix/store/908kic4k4zb6wvnw10p5q1363ym94qgj-perl-5.30.1/lib/perl5/5.30.1
17:14 -!- Irssi: process 0 (perl -V) terminated with return code 0
```

That's a long list of library paths. The [higlighter script][highlighter]
is now able to load properly.

## References
- <https://github.com/irssi/irssi/blob/master/INSTALL>
- <https://stackoverflow.com/q/841785/90674>
- <https://perlmaven.com/how-to-change-inc-to-find-perl-modules-in-non-standard-locations>
- <https://github.com/NixOS/nixpkgs/blob/master/pkgs/top-level/perl-packages.nix>
- <https://github.com/NixOS/nixpkgs/blob/1c279a00119beef0e09a29b844296dc829ca9e2d/pkgs/servers/monitoring/munin/default.nix#L120>
- <https://github.com/NixOS/nixpkgs/issues/10350>

[scriptassist.pl]: https://github.com/irssi/irssi/blob/master/scripts/scriptassist.pl
[highlighter]: https://github.com/stigok/notification-hub/blob/master/examples/irssi/hilightcmd.pl
[notification-hub]: https://github.com/stigok/notification-hub/
