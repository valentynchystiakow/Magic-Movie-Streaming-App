// imports libraries(components)
import axios from "axios";

// defines api base url
const apiUrl = import.meta.env.VITE_API_BASE_URL;

// creates and exports axios instance with properties
export default axios.create({
  baseURL: apiUrl,
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});
