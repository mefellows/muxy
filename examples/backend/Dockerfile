FROM gliderlabs/alpine
MAINTAINER Matt Fellows <matt.fellows@onegeek.com.au>

# For the lazy, run the gox command locally
# uncomment this

# ADD microgo microgo
# ENV PORT 80
# EXPOSE 80
# ENTRYPOINT ["/microgo"]

# For the build server - longer but automated
FROM golang:1.5.1-onbuild
ENV PORT 80
EXPOSE 80