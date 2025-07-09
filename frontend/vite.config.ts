import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@common": path.resolve(__dirname, "./src/common"),
      "@lib": path.resolve(__dirname, "./src/lib"),
      "@store": path.resolve(__dirname, "./src/store"),
    },
  },
});
