# Makefile

dev:
	npx concurrently \
		"npx @tailwindcss/cli -i ./static/css/style.css -o ./static/generated/tailwind.css --verbose --watch" \
		"air"

