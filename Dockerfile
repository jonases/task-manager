FROM golang:1.8

EXPOSE 8001

#RUN mkdir /go/src/nciproject
ADD . /go/src/github.com/jonases/cybersecuryproject

#ENV GOPATH /go/src/github.com/jonases/cybersecuryproject
#ENV GOBIN /go/src/github.com/jonases/cybersecuryproject
#ENV PATH $PATH:$GOPATH

#RUN go get "github.com/IBM-Bluemix/go-cloudant"
#RUN go get "github.com/timjacobi/go-couchdb"
#RUN go get "golang.org/x/crypto/bcrypt"
#RUN go get "github.com/josephspurrier/csrfbanana"
#RUN go get "github.com/gorilla/sessions"
#RUN go get "github.com/gorilla/mux"

ENV CLOUDANT_USER_NAME <CLOUDANT_USERNAME>
ENV CLOUDANT_PASSWORD <CLOUDANT_PASSWORD>

#RUN go install -v -gcflags "-N -l" /go/src/github.com/jonases/cybersecuryproject
RUN cd /go/src/github.com/jonases/cybersecuryproject && make
#RUN cd /go/src/nciproject/src && go install -v -gcflags "-N -l" main.go
#RUN cd /go/src/nciproject/src && go build -o webapp

ENTRYPOINT [ "/go/src/github.com/jonases/cybersecuryproject/cybersecuryproject" ]
