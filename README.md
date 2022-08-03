Goonware is still very much a work-in-progress and you probably won't find much
use in it as an end-user right now. Come back later :3

# Goonware

Edgeware but in Go for cumming performance™

![](.github/ui.png)

# Features

- Complete feature parity with Edgeware, including Edgeware package support
- Completely cross-platform
- Portable compiled binary = much simpler to set up (and it's faster I guess)
- More intuitive UI with tooltips
- A new collection of "Passive" features

# Build Instructions

Just `go build`; on Linux you might need to install `libx11`, and on Windows you need GCC --
the it just worksTM solution is https://jmeubank.github.io/tdm-gcc/.

The first build will take a long while. Subsequent builds will be a lot faster, but still slow
by Go standards because the GUI library depends on a lot of C.