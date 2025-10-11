# Makefile

dev:
	npx concurrently \
		"npx @tailwindcss/cli -i ./static/css/style.css -o ./static/css/tailwind.css --verbose --watch" \
		"air"

