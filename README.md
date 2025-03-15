# `Pretty logger` and `Spell checker` on Go

This is just a small project that at first was about spell checking, but at some points I wanted to check out what a slog package could do.

It wasn't quite right, so I decided to write my own logger package and then finsih the spell checking algorithm.

### `Pretty logger` package

Contain a customly written logger with methods:
- Debug/Info/Warn/Error/Fatal - for printing logs;
- Set methods - to change output information and format;
- Pretty output - all logs are formated strings with additional info of caller function, current date-time and level of log.

### `Spell checker` algorithm

This is 'Damerau–Levenshtein distance' that in O(n^2) time finds a simillarity score of 2 strings. Then whole spell checker just iterates through a whole dictionary to find best matching words.

Spell checker separeted into package for a better testing purposes.
