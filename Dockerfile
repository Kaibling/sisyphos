FROM golang:1.19
#RUN apt update && 
ENTRYPOINT ["tail", "-f", "/dev/null"]