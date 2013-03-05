# whac-a-gopher: Simple game written in Go

This is an example project for my speak in [devfestW Korea 2013][3]
which is about Go's power of handling concurrency.

![ScreenshotOfTheGame](https://lh4.googleusercontent.com/-QqMLiacnqaE/UTRT047USNI/AAAAAAAACSA/52djbhOWATI/s625/%EC%8A%A4%ED%81%AC%EB%A6%B0%EC%83%B7%2C+2013-03-04+16%3A55%3A08.png)

[Video(youtube)](http://youtu.be/fqvJWG4cWIg)

## How to build

[Install Go][1]

### Install requirement packages for Go-SDL (on Ubuntu 12.04)

Install dependency packages:

    $ sudo apt-get install libsdl1.2-dev libsdl-mixer* libsdl-image* libsdl-ttf*

On under Ubuntu 12.10, need to make `SDL_ttf.pc` to `/usr/lib/pkgconfig` with
following context:

    prefix=@prefix@
    exec_prefix=@exec_prefix@
    libdir=@libdir@
    includedir=@includedir@

    Name: SDL_ttf
    Description: ttf library for Simple DirectMedia Layer with FreeType 2 support
    Version: @VERSION@
    Requires: sdl >= @SDL_VERSION@
    Libs: -L${libdir} -lSDL_ttf
    Cflags: -I${includedir}/SDL

[Read more][2] about this issue.

### Build and Install

    $ go get
    $ go build

## Usage

    $ ./whac-a-gopher

[1]:http://golang.org/doc/install
[2]:https://github.com/banthar/Go-SDL/issues/35#issuecomment-3597261
[3]:https://sites.google.com/site/2013devfestwkorea
