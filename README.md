# `Spell checker`

### Overview:

This is just a small project that at first was about spell checking, but at some points I wanted to check out what a slog package could do.

It wasn't quite right, so I decided to write my own logger package and then finsih the spell checking algorithm.

[Customly written logger now splitted into new repository.](https://github.com/TrueHopolok/plog)

---

### Spell checker algorithm

This is [Damerau–Levenshtein distance](https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance) that in O(n^2) time finds a simillarity score of 2 strings. Then whole spell checker just iterates through a whole dictionary to find best matching words.

Spell checker separeted into package for a better testing purposes.
