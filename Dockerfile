FROM public.ecr.aws/docker/library/alpine:3.20

WORKDIR /
ADD flowstatesrv /flowstatesrv

CMD ["/flowstatesrv"]
