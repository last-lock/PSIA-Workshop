# Use a base image with Java installed
FROM openjdk:11-jre-slim

# Set environment variables
ENV KAFKA_VERSION=3.4.0
ENV SCALA_VERSION=2.13
ENV KAFKA_HOME=/opt/kafka

# Install necessary tools
RUN apt-get update && apt-get install -y wget tar

# Download and extract Kafka
RUN wget -q https://downloads.apache.org/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -O /tmp/kafka.tgz
RUN tar xzf /tmp/kafka.tgz -C /opt 
RUN ln -s /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION} ${KAFKA_HOME} 
RUN  rm /tmp/kafka.tgz

# Set up configuration
COPY server.properties ${KAFKA_HOME}/config/
COPY zookeeper.properties ${KAFKA_HOME}/config/

# Expose Kafka and ZooKeeper ports
EXPOSE 9092 2181

# Set up volumes for Kafka and ZooKeeper data
VOLUME ["/kafka-logs", "/zookeeper-data"]

# Start ZooKeeper and Kafka
CMD ${KAFKA_HOME}/bin/zookeeper-server-start.sh ${KAFKA_HOME}/config/zookeeper.properties & \
    ${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.properties