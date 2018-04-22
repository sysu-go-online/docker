FROM docker
ADD main /

ENTRYPOINT ["/main"]