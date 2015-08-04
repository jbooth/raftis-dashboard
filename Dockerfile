FROM gliderlabs/alpine:3.1
ADD dashboard /bin/dashboard
EXPOSE 8080
CMD /bin/dashboard
