FROM yindaheng98/go-git
WORKDIR /app
COPY . .
COPY ./id_rsa /root/.ssh/id_rsa
COPY ./id_rsa.pub /root/.ssh/id_rsa.pub

RUN git config --global url."git@gitee.com:".insteadOf "https://gitee.com/" && \
    chmod 0600 /root/.ssh/id_rsa && \
    echo "StrictHostKeyChecking no " > /root/.ssh/config && \
    go get -d -v ./... && \
    go build main -v -o /PressureMeterNet-Master

FROM egaillardon/jmeter
STOPSIGNAL SIGINT

RUN mkdir /PressureMeterNet && \
    mkdir /PressureMeterNet/Data && \
    chmod a+rw /PressureMeterNet/Data
COPY --from=0 /PressureMeterNet-Master /PressureMeterNet
WORKDIR /PressureMeterNet

EXPOSE 8080
VOLUME [ "/PressureMeterNet/Data" ]
CMD ["/PressureMeterNet/PressureMeterNet-Master"]