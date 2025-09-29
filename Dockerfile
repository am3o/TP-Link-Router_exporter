FROM scratch

COPY TP-Link-Router_exporter /TP-Link-Router_exporter

ENTRYPOINT ["/TP-Link-Router_exporter"]
