import axios from "axios";

const BASE_API_URL = import.meta.env.VITE_BASE_API_URL;
const BASE_API_TOKEN = import.meta.env.VITE_BASE_API_TOKEN;

// console.log(`base url: ${BASE_API_URL}`);
// console.log(`token : ${BASE_API_TOKEN}`);

const http = axios.create({
  baseURL: `${BASE_API_URL}/api`,
  headers: {
    "Content-Type": "application/json",
    Authorization: BASE_API_TOKEN,
  },
});

export default http;
