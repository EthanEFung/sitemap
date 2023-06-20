# Sitemap

Simple command line program that generates and writes a domain's xml sitemap to
standard output. Given a domain, the program makes an http GET request expecting
text/html. Upon response, the program finds all the links on the page and repeatedly
makes http requests for every link until all of the domain's links are found in a
breadth first search. Note, this does not guarentee every url is reached for a
domain, only every url reachable from visiting the original domain url.

Note: This program has *NOT* been thoroughly tested, and should not be used for
production! This program was written as a learning exercise, and should be treated
as such.

## How to use

Be sure to have a modern version of Golang. This program was written with `1.19`.
In your terminal clone the repo, install the dependencies, and build

```bash
git clone git@github.com:EthanEFung/sitemap.git
cd sitemap
go install
sitemap -url http://www.someblog.com
```

You can even write the map to a file. On linux or macOS, you can pipe the output
to a file with `tee`

```bash
sitemap -url http://www.someblog.com | tee ~/Downloads/someblog_sitemap.xml
```

## Caveats

The underlying algorithm is rudimentary and runs in linear time and space. Meaning
that creating a sitemap for large websites is not viable in its current state. 

The algorithm also is not fully supporting relative paths. The program currently
only supports absolute and root directory hrefs. I plan at a later time to address
this issue.

Again, this project was built for educational purposes and this repo will be mostly
unmaintained, but I do welcome pull requests.

## Interesting problems to solve

1. As noted in the "Caveats" section. `./` and `../` are not supported. Some effort can be placed into this. 
2. More interestingly would be solving for potential memory issues caused by generating sitemaps on large websites. In terms of performance, some optimizations to explore are: utilize a [flyweight design pattern](https://refactoring.guru/design-patterns/flyweight) for the queue instead of a linked list, and utilize a disc instead of RAM to store the seen urls.
3. The program is currently only using one routine to generate the sitemap. What would it look like to use go routines to visit all the of the endpoints of a domain?
4. Currently the output is spitting out xml. Why not dumb it down to just create a text file of the links? That way we can use the text file to generate xml, json, or any other format instead.
5. Add a depth flag so that given a certain point, the breadth first search will stop. 
6. Having a depth flag, provide a way to merge two sitemaps into one, making it easier to break the task of creating a large websites sitemap into chunks.
7. Add logging
