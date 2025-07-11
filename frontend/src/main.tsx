import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import ToastProvider from "@/pages/common/ToastProvider.tsx";

import "./global.css";

import App from "./App.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
    <ToastProvider />
  </StrictMode>
);
