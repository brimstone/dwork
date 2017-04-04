FROM scratch

COPY app /

ENTRYPOINT ["/dwork"]
