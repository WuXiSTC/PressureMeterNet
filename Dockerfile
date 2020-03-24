FROM yindaheng98/go-git
WORKDIR /app
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

COPY --from=0 /PressureMeterNet-Master /jmeter
ADD Config.yaml /jmeter

EXPOSE 8080
WORKDIR /jmeter
RUN mkdir Data
VOLUME [ "/jmeter/Data" ]
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/jmeter/PressureMeterNet-Master"]