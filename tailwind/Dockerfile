FROM node:22
WORKDIR /app
RUN npm install tailwindcss
CMD [ "npx", "tailwindcss", "-i", "./tailwind/input.css", "-o", "./public/tailwind.css", "--watch" ]