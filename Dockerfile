
# Use the official MySQL image
FROM mysql:latest

# Set environment variables for MySQL root password
ENV MYSQL_ROOT_PASSWORD=mypassword

# Expose the MySQL port (3306)
EXPOSE 3306

# Create a volume for persistent data
VOLUME ["/var/lib/mysql"]

# Copy any initialization scripts to be executed on startup
COPY ./init.sql /docker-entrypoint-initdb.d/

# Run initialization scripts on container startup
CMD ["mysqld"]
