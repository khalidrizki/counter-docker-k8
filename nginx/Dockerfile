FROM nginx:alpine

COPY frontend /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /var/cache/nginx \
    /var/run/nginx \
    /var/log/nginx \
    && chgrp -R 0 /var/cache/nginx /var/run/nginx /var/log/nginx /usr/share/nginx/html \
    && chmod -R g=u /var/cache/nginx /var/run/nginx /var/log/nginx /usr/share/nginx/html