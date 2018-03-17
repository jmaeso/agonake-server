FROM alpine:3.7

COPY ./bin/agonake-server /home/agonake/server
RUN adduser -D server && \
    chown -R server /home/agonake && \
    chmod o+x /home/agonake/server

USER server
ENTRYPOINT /home/agonake/server
