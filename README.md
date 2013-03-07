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

On under Ubuntu 12.10, you might make some `.pc`s manunually.
[Read more][2] about this issue.

### Download, Build and Install

    $ go get github.com/suapapa/whac-a-gopher

> The executable is created in $GOPATH/bin/

## Run

Run in default settings, 9 gophers in window mode.

    $ whac-a-gopher

Try various combinations of options. For example, 45 gophers in FHD.

    $ whac-a-gopher -w 1920 -h 1080 -f true

And refer help, `-help` for more.


[1]:http://golang.org/doc/install
[2]:https://github.com/banthar/Go-SDL/issues/35#issuecomment-3597261
[3]:https://sites.google.com/site/2013devfestwkorea
