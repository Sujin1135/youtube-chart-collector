FROM ubuntu:latest
LABEL authors="choemingyu"

ENTRYPOINT ["top", "-b"]