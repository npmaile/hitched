import { defineConfig } from "vite"
import { imagetools } from "vite-imagetools"
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
	plugins: [
		imagetools(),
		react(),
		tailwindcss()
	]
})
