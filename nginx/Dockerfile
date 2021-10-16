FROM nginx:alpine

RUN rm -v /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/nginx.conf
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

CMD ["/usr/sbin/nginx"]
