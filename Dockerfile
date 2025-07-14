#    docker build -t imposter-img .
#    docker run -d --name imposter-cont -p 3306:3306 imposter-img                                                                                                   ─╯
# Use the official MySQL image
FROM mysql:8.0.35

# Set environment variables for MySQL root password
ENV MYSQL_ROOT_PASSWORD=mypassword

# Copy the MySQL configuration file into the container
COPY my.cnf /etc/mysql/conf.d/

# Expose the MySQL port (3306)
EXPOSE 3306

# Create a volume for persistent data
VOLUME ["/var/lib/mysql"]

# Run initialization scripts on container startup
CMD ["mysqld"]
