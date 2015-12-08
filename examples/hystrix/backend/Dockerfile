FROM gliderlabs/alpine
#FROM golang:1.5.1-onbuild
MAINTAINER Matt Fellows <matt.fellows@onegeek.com.au>

# gox -osarch="linux/amd64" -output="microgo"
# Remove this for on build
ADD microgo microgo
ENV PORT 80
EXPOSE 80
# Remove this for onbuild
ENTRYPOINT ["/microgo"]
