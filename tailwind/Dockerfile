FROM node:22
WORKDIR /app
RUN npm install tailwindcss@3.4.10
CMD [ "npx", "--yes", "tailwindcss", "-c", "./tailwind/tailwind.config.js", "-i", "./tailwind/input.css", "-o", "./public/tailwind.css", "--watch" ]