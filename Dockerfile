FROM scratch
COPY lucky /
EXPOSE 16601
# Store config and data in /goodluck
WORKDIR /goodluck
ENTRYPOINT ["/lucky"]
CMD ["-c", "/goodluck/lucky.conf"]
