_default:
	@just --list

setup: tailwind tailwind-prod

tailwind:
	tailwindcss -i app/public/in.css -o app/public/static/tailwind.css --config app/tailwind.config.js

tailwind-prod:
	tailwindcss -i app/public/in.css -o app/public/static/tailwind.css --config app/tailwind.config.js
