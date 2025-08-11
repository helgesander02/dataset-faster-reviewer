# dockerfile for building and running the frontend application
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# final stage to run the application
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app ./
CMD ["npm", "run", "start"]
EXPOSE 3000