---
layout: post
title:  "Numbering the filenames when downloading multiple files with cURL"
date:   2021-07-14 10:43:43 +0200
categories: linux curl
excerpt: When downloading multiple files with curl I want the filenames starting with its number in line of files.
#proccessors: pymd
---

When downloading a list of URL's using `curl -O <source>` the resulting files
will get the same filename as they appear in the source (`basename <source>`).

I am downloading a list of MP3s where the filenames only contain the names of
the guests. This means that the files aren't sorted properly when they are
saved to my local disk. I want to number the files as they appear.

```bash
# Start counting from 1
i=1

# I am reading a HTML page containing a list of podcasts starting with the newest one,
# so the last item is the first episode.
# Get all links ending with '.mp3' and reverse the list with tac.
urls=$(grep -oP '[^"]+\.mp3' src.html | tac)

for url in $urls
do
  # Output the contents of the URL to e.g. "1-the-first-episode.mp3"
  curl -o "${i}-$(basename "$url")" -L "$url"

  # Increment the counter
  i=$(( i + 1 ))
done
```

This gives me a list like this

```
10-eldar_vagan.mp3
11-jon_niklas_ronning.mp3
12-tom_mathisen.mp3
13-stephen_ackles.mp3
14-marianne_krogness.mp3
15-lars_eikanger.mp3
16-guri_schanke.mp3
17-anders_hatlo.mp3
18-brit_elisabeth_haagensli.mp3
19-terje_formoe.mp3
1-bak-scenen-john-nyutstumo.mp3
1-tore_ryen.mp3
20-espen_beranek_holm.mp3
21-halvdan_sivertsen.mp3
22-hilde_lyran.mp3
23-arne_garvang.mp3
24-anne_marie_ottersen.mp3
25-dag_vagsas.mp3
26-karl_sundby.mp3
27-finn_kalvik.mp3
28-benny_borg.mp3
29-endre_lundgren.mp3
2-trine_rein.mp3
30-niklas_baarli.mp3
31-else_kass_furuseth.mp3
32-hans_morten_hansen.mp3
33-joakim_radiomann.mp3
34-tom_sterri.mp3
35-arve_sigvaldsen.mp3
36-age_sten_nilsen.mp3
37-harald_maele.mp3
38-steinar_albrigtsen.mp3
39-solfrid_heier.mp3
3-william_kristoffersen_ole_ivars.mp3
40-anita_skorgan.mp3
41-rune_rudberg.mp3
42-nils_vogt.mp3
43-sturla_berg_johansen.mp3
44-hallvard_flatland.mp3
45-viggo_sandvik.mp3
46-anitra_eriksen.mp3
47-henriette_lien.mp3
48-1-bak-scenen-john-nyutstumo.mp3
49-marit-synnove-berg_ferdig.mp3
4-hellbillies.mp3
5-havard_bakke.mp3
6-christian_ingebrigtsen.mp3
7-tommy_michaelsen_picazzo.mp3
8-per_oystein_sorensen_fra_lippo_lippi.mp3
9-per_christian_ellefsen.mp3
```

Ideally I'd like to zero-pad the numbers so they get sorted properly.
I can use `printf` for that.

```bash

for url in $urls
do
  # Save the zeropadded number to a variable named "zeropadded".
  # The number 02 says you want to pad with 0, to make the string 2 characters long.
  printf -v zeropadded "%02d" $i

  # Output the contents of the URL to e.g. "1-the-first-episode.mp3"
  curl -o "${zeropadded}-$(basename "$url")" -L "$url"

  # Increment the counter
  i=$(( i + 1 ))
done
```

Great! Now the first 10 items in the list looks like this

```
01-tore_ryen.mp3
02-trine_rein.mp3
03-william_kristoffersen_ole_ivars.mp3
04-hellbillies.mp3
05-havard_bakke.mp3
06-christian_ingebrigtsen.mp3
07-tommy_michaelsen_picazzo.mp3
08-per_oystein_sorensen_fra_lippo_lippi.mp3
09-per_christian_ellefsen.mp3
10-eldar_vagan.mp3
11-jon_niklas_ronning.mp3
```

## References
- <https://stackoverflow.com/a/8789815/90674>
