FROM golang:1.20.0
RUN apt-get update && apt-get install -y unzip pkg-config libaio1 curl

RUN cd /tmp && \
    curl -o instantclient-basiclite.zip https://download.oracle.com/otn_software/linux/instantclient/instantclient-basiclite-linuxx64.zip -SL && \
    unzip instantclient-basiclite.zip && \
    mv instantclient*/ /usr/lib/instantclient && \
    rm instantclient-basiclite.zip && \
    ln -s /usr/lib/instantclient/libclntsh.so.19.1 /usr/lib/libclntsh.so && \
    ln -s /usr/lib/instantclient/libocci.so.19.1 /usr/lib/libocci.so && \
    ln -s /usr/lib/instantclient/libociicus.so /usr/lib/libociicus.so && \
    ln -s /usr/lib/instantclient/libnnz19.so /usr/lib/libnnz19.so && \
    ln -s /usr/lib/libnsl.so.2 /usr/lib/libnsl.so.1 && \
    ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 && \
    ln -s /lib64/ld-linux-x86-64.so.2 /usr/lib/ld-linux-x86-64.so.2

ENV ORACLE_BASE=/usr/lib/instantclient
ENV LD_LIBRARY_PATH=/usr/lib/instantclient
ENV TNS_ADMIN=/usr/lib/instantclient
ENV ORACLE_HOME=/usr/lib/instantclient

RUN go install github.com/cespare/reflex@latest

#COPY reflex.conf /usr/local/etc/
#COPY build.sh /usr/local/bin/

WORKDIR /app

VOLUME /go

#USER root 
#RUN chmod 755 /usr/local/bin/build.sh

# CMD ["reflex", "-d", "none", "-c", "/usr/local/etc/reflex.conf"]