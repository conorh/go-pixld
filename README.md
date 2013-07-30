== go-pixeld

An image resizing server written in Go.

Note: still in development, but basic functionality working.

TODO:
* Decode other image formats - png using image/png
* Use imagemagick convert for speed
* Tests for basic functionality
* More conversion types

== Usage

1. Run the server. 

    ./go-pixeld

2. Request an image via the server:

    http://server:8080/s?img=http://someurl.com/someimage.jpg&w=100&h=200

The requested image will be downloaded, resized, cached and served up. Future requests
for the same image will be loaded from the file cache.

== License

The MIT License (MIT)

Copyright (c) 2013 Conor Hunt <conor.hunt@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
