# Use an image containing environment for building Go programs.
FROM golang:1.15 AS build

# Set the working directory inside the container.
WORKDIR /src

# Donwload all dependencies.
COPY go.* ./
RUN go mod download

# Build the source code.
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o kvstore

# Scratch image contains no distribution files. The resulting image and
# containers will have only our statically linked application binary.
FROM scratch AS image

COPY --from=build /src/kvstore /bin/kvstore

# Tell Docker we'll be using port 8000.
EXPOSE 8000

# Tell Docker to execute this command on a "docker run".
CMD ["/bin/kvstore"]
