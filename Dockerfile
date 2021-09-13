FROM golang:1.17

WORKDIR /app
COPY . /app/

RUN make

CMD ./golox
