# syntax=docker/dockerfile:1

FROM golang:1.19 AS build-stage

#######################################################
# Copy the go.mod and go.sum file into your project 
# directory /app which, owing to your use of WORKDIR, 
# is the current directory (./) inside the image. 
# Unlike some modern shells that appear to be 
# indifferent to the use of trailing slash (/), and can
# figure out what the user meant (most of the time), 
# Docker's COPY command is quite sensitive in its 
# interpretation of the trailing slash.
#######################################################
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./

#######################################################
# Now that you have the module files inside the Docker 
# image that you are building, you can use the RUN 
# command to run the command go mod download there as 
# well. This works exactly the same as if you were 
# running go locally on your machine, but this time 
# these Go modules will be installed into a directory 
# inside the image.
#######################################################
RUN go mod download && go mod verify

#######################################################
# Copy the source into this container
#######################################################
COPY . .

#######################################################
# Compile your application
# The result of that command will be a static 
# application binary named docker-gs-ping and located 
# in the root of the filesystem of the image that you 
# are building. You could have put the binary into any 
# other place you desire inside that image, the root 
# directory has no special meaning in this regard. It's 
# just convenient to use it to keep the file paths 
# short for improved readability.
#######################################################
RUN CGO_ENABLED=0 go build -v -o ./fsm .

# Deploy the application binary into a lean image
FROM alpine AS release-authentication

RUN mkdir /app
WORKDIR /app

COPY --from=build-stage /app .

RUN chmod +x /app/fsm

ENTRYPOINT ["./fsm"]