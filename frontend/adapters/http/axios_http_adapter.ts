import { HttpPort, RequestConfig } from "@/core/ports/http/http_port";
import axios, { AxiosInstance } from "axios";

export class AxiosHttpAdapter implements HttpPort {
    private client: AxiosInstance;
    constructor(baseURL: string) {
        this.client = axios.create({
            baseURL,
            // timeout: 5000, //5 seconds
            withCredentials: true, // cookies
            headers: {
                "Content-Type": "application/json",
            },
        });
        // 🔥 Attach token dynamically
        // this.client.interceptors.request.use(config => {
        //     const token = token.getAccessToken();
        //     if (token) {
        //         config.headers.Authorization = `Bearer ${token}`;
        //     }
        //     return config;
        // });

        // // 🔥 production-ready interceptor
        // this.client.interceptors.response.use(
        //     res => res,
        //     err => {
        //         // Log only the status and data to avoid the massive red stack trace in the console
        //         console.warn(`API Error [${err.response?.status}]:`, err.response?.data);

        //         const mapAxiosError = MapAxiosError(err)
        //         // if (mapAxiosError.statusCode == 401) {
        //         //     const currectParams = window.location.search;
        //         //     window.location.href = `/${API_ENDPOINTS.Auth.signin + currectParams}`;

        //         // }
        //         throw mapAxiosError;
        //     }
        // );
    }

    async get<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.get<T>(config.url, { params: config.params });
        return res.data;
    }
    async post<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.post<T>(
            config.url,
            config.body,
            { params: config.params }
        );
        return res.data;
        // return await this.client.post<T>(config.url, config.body);
    }
    async put<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.put<T>(
            config.url,
            config.body,
            { params: config.params }
        );
        return res.data;
    }
    async delete<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.delete<T>(config.url, { params: config.params });
        return res.data;
    }
}