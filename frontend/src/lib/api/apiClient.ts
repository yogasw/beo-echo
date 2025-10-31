import { BASE_URL_API } from "$lib/utils/authUtils";
import axios from "axios";

// Create axios instance with default config
const apiClient = axios.create({
    baseURL: BASE_URL_API,
    headers: {
        'Content-Type': 'application/json'
    }
});

export { apiClient };
export default apiClient;
