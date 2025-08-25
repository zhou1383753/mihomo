FROM metacubex/mihomo:Alpha AS m
FROM alpine AS upxer
RUN apk add upx
ARG arch
COPY --from=m /mihomo /usr/share/bin/mihomo
COPY ./build/clash-admin_linux_${arch} /usr/share/bin/linux
RUN upx -9 /usr/share/bin/linux
RUN upx -9 /usr/share/bin/mihomo

FROM metacubex/mihomo:Alpha AS mihomo
RUN rm -rf /mihomo
COPY --from=upxer /usr/share/bin/mihomo /mihomo
RUN apk add tini
ARG arch
RUN echo $arch
RUN if [ "$arch" = "arm" ];then rm -r /root/.config/mihomo/geoip.dat /root/.config/mihomo/geosite.dat; fi
RUN if [ "$arch" = "arm64" ];then rm -r /root/.config/mihomo/geoip.dat /root/.config/mihomo/geosite.dat; fi

FROM scratch
COPY --from=upxer /usr/share/bin/linux /usr/share/bin/linux
COPY --from=mihomo / / 


COPY config.yaml /root/.config/mihomo/config.yaml
COPY clash_ui /root/.config/mihomo/ui
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN ln -sf /root/.config/mihomo /etc/clash
RUN ln -sf /root/.config/mihomo /etc/mihomo
COPY clash_start.sh /etc/clash/start.sh
RUN touch /var/log/run.log

COPY timezone_cst /etc/timezone
COPY localtime_shanghai /etc/localtime
ENV TZ=Asia/Shanghai
WORKDIR /root/clash-admin
EXPOSE 8080
ENTRYPOINT [ "/docker-entrypoint.sh" ]