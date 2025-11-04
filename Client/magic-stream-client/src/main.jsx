// imports libraries(components)
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.jsx";
// imports bootstrap classes
import "bootstrap/dist/css/bootstrap.min.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthProvider.jsx";

// finds the root element in the index.html file and renders the App component
createRoot(document.getElementById("root")).render(
  <StrictMode>
    {/* Wraps the app in the auth provider component for authentication*/}
    <AuthProvider>
      {/* Wraps the app in a browser router for routing between pages*/}
      <BrowserRouter>
        <Routes>
          <Route path="/*" element={<App />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  </StrictMode>
);
