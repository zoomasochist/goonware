# Goonware

Edgeware but in Go for cumming performance™

![](.github/ui.png)

# Build Instructions

Just `go build`; on Linux you might need to install `libx11`, and on Windows you need GCC --
the it just worksTM solution is https://jmeubank.github.io/tdm-gcc/.

The first build will take a long while. Subsequent builds will be a lot faster, but still slow
by Go standards because the GUI library depends on a lot of C.