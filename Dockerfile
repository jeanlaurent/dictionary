FROM scratch

ADD app /app
COPY dictionary /

EXPOSE 8080

CMD ["/dictionary"]