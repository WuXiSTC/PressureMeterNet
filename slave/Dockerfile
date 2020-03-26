FROM yindaheng98/go-git
COPY . .
RUN go get -d -v ./... && \
    go build -v -o /PressureMeterNet-Slave

FROM egaillardon/jmeter
STOPSIGNAL SIGINT

COPY --from=0 /PressureMeterNet-Slave /jmeter
COPY entrypoint.sh /usr/local/bin/
WORKDIR /jmeter

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/jmeter/PressureMeterNet-Slave", "-h"]