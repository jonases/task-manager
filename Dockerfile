FROM alpine:3.7

EXPOSE 8001

ADD task-manager-linux64 /task-manager
RUN mkdir /public
ADD public /public

RUN apk --no-cache add ca-certificates
COPY cloudantcert.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

ENV CLOUDANT_USER_NAME <username>
ENV CLOUDANT_PASSWORD <password>

ENTRYPOINT [ "./task-manager" ]
