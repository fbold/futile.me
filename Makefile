# Makefile

dev:
	npx @tailwindcss/cli --watch -i ./static/css/style.css -o ./static/css/tailwind.css &\
	air

