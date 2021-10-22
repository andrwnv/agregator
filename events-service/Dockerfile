FROM node:latest

WORKDIR /src/app

COPY package.json /src/app/package.json
RUN npm install

COPY . /src/app

EXPOSE 3020
