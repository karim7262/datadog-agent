FROM multiarch/crossbuild
RUN apt-get update
RUN apt-get install -y rake curl
RUN curl -sL -o /bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
RUN chmod +x /bin/gimme
ADD . /go/src/github.com/DataDog/datadog-agent
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH
ENV CROSS_TRIPLE=x86_64-apple-darwin
env GOOS=darwin
env GOARCH=x86_64
WORKDIR /go/src/github.com/DataDog/datadog-agent
