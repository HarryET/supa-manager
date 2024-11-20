FROM node:16

WORKDIR /app

COPY studio/ .
COPY studio/patch.sh /app/patch.sh
COPY studio/build.sh /app/build.sh

RUN chmod +x /app/patch.sh /app/build.sh
RUN ./patch.sh && ./build.sh

EXPOSE 3000
CMD ["npm", "start"] 