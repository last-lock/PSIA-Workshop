FROM alpine:3.19
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.license=GPL-3.0 \
    org.label-schema.name=samba-dc \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url=https://github.com/instantlinux/docker-tools

ENV ADMIN_PASSWORD_SECRET=samba-admin-password \
    ALLOW_DNS_UPDATES=secure \
    BIND_INTERFACES_ONLY=yes \
    DOMAIN_ACTION=provision \
    DOMAIN_LOGONS=yes \
    DOMAIN_MASTER=no \
    INTERFACES="lo eth0" \
    LOG_LEVEL=1 \
    MODEL=standard \
    NETBIOS_NAME=dc\
    REALM=ad.example.com \
    SERVER_STRING="Samba Domain Controller" \
    TZ=CST \
    WINBIND_USE_DEFAULT_DOMAIN=yes \
    WORKGROUP=AD \
    DNS_PORT=5353 
    # DNS_PORT smb.conf.j2 modification


ARG SAMBA_VERSION=4.18.9-r0

COPY *.conf.j2 /root/
COPY entrypoint.sh /usr/local/bin/
RUN apk add --update --no-cache krb5 ldb-tools samba-dc=$SAMBA_VERSION samba-winbind-clients=$SAMBA_VERSION tdb \
      bind bind-libs bind-tools libcrypto3 libxml2 tzdata py3-setuptools py3-pip && \
    pip install j2cli --break-system-packages && \
    apk del py3-pip && \
    chmod 0755 /usr/local/bin/entrypoint.sh

#FIXME Update user password
RUN adduser -D testUser && \
    echo "testUser:password1!" | chpasswd
RUN (echo "password1!"; echo "password1!") | smbpasswd -s -a testUser


VOLUME /etc/samba /var/lib/samba
EXPOSE 53 53/udp 88 88/udp 135 137-138/udp 139 389 389/udp 445 464 464/udp 636 3268-3269 49152-65535

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]