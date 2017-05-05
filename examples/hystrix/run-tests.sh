#!/bin/bash

function clean {
  docker-compose stop
  docker-compose rm -f
}

function build {
  cd api; gox -osarch="linux/amd64" -output="microgo"; cd ..
  cd backend; gox -osarch="linux/amd64" -output="microgo"; cd ..
  docker-compose build
}

function run {
  docker-compose up
}

case $1 in
"run")
  run
  ;;
"build")
  build
  ;;
"clean")
  clean
  ;;
*)
  printf "Cleaning, building and running!\n"
  clean
  build
  run
  ;;
esac
