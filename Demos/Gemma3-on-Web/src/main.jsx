/**
 * Entry point of the React application.
 *
 * This file sets up the root of the React application and renders the main App component
 * within a StrictMode wrapper for highlighting potential problems in the application.
 */
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Analytics } from "@vercel/analytics/react";
import { SpeedInsights } from "@vercel/speed-insights/react";
import "./styles/index.css";
import App from "./App.jsx";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <SpeedInsights />
    <Analytics />
    <App />
  </StrictMode>
);
