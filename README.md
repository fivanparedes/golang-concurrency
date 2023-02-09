# golang-concurrency
Implementation of a comic downloader written in Go 1.20, using work and pooling patterns. This is a demonstration of concurrent programming for the Paradigms &amp; Programming Languages course of the Universidad Nacional de Misiones.

## Installation
It is not necessary to have Go runtime and compiler in the machine. Just use Docker!
* Download the certificate of xkcd.com and put it in the `certs` directory.
* Build the image with `docker build . -t golang-concurrency`
* Run the container with `docker run golang-concurrency`.


