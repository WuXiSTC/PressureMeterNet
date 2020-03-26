FROM yindaheng98/go-git
COPY . .
COPY ./id_rsa /root/.ssh/id_rsa
COPY ./id_rsa.pub /root/.ssh/id_rsa.pub

RUN git config --global url."git@gitee.com:".insteadOf "https://gitee.com/" && \
    chmod 0600 /root/.ssh/id_rsa && \
    echo "StrictHostKeyChecking no " > /root/.ssh/config && \
    go get -d -v ./... && \
    go build -v -o /PressureMeterNet-Master

FROM egaillardon/jmeter
STOPSIGNAL SIGINT

RUN mkdir /jmeter/Data && chmod a+rw /jmeter/Data
COPY --from=0 /PressureMeterNet-Master /jmeter
COPY entrypoint.sh /usr/local/bin/
WORKDIR /jmeter

EXPOSE 8080
VOLUME [ "/jmeter/Data" ]
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/jmeter/PressureMeterNet-Master", "-h"]