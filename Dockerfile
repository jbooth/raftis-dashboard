FROM debian:jessie
ADD dashboard /bin/dashboard
EXPOSE 8080
CMD /bin/dashboard
