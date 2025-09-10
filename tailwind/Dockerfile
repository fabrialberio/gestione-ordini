FROM node:22
WORKDIR /app
RUN npm install tailwindcss@3.4.10 @tailwindcss/cli
CMD [ "npx", "--yes", "@tailwindcss/cli", "-c", "./tailwind/tailwind.config.js", "-i", "./tailwind/input.css", "-o", "./public/tailwind.css", "--watch" ]
