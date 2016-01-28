FROM iron/base

ADD sshttp /sshttp/sshttp

ENTRYPOINT ["/sshttp/sshttp"]
