FROM alpine:3.19
WORKDIR /app
RUN apk --no-cache upgrade
COPY urlshortener .
EXPOSE 9808
CMD ["./urlshortener"]