FROM nunchistudio/blacksmith-enterprise:0.12.0-alpine

# Please do not mind the following lines as they are just here for our team
# when working in local with a Go modules proxy.
# START UNCOMMENT(proxy)
# ENV GOPROXY=http://host.docker.internal:3000
# ENV GONOSUMDB=github.com/nunchistudio/*
# END UNCOMMENT(proxy)

ADD ./ /smithy
WORKDIR /smithy

RUN rm -rf go.sum
RUN go mod tidy
