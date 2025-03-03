FROM nginx:1.24

COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /var/log/nginx && \
    chown -R nginx:nginx /var/log/nginx 