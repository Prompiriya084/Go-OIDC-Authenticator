import { HttpPort, RequestConfig } from "@/core/ports/http/http_port";
import axios, { AxiosError, AxiosInstance } from "axios";
import { MapAxiosError } from "./map_axios_error";

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
        // 🚀 Production-Ready Response Interceptor
        this.client.interceptors.response.use(
            (res) => res,
            (err: AxiosError) => {
                // Log รายละเอียดความพังบน Terminal ฝั่ง Server โดยไม่พ่น Stack Trace สีแดงเต็มหน้าจอ
                console.warn(
                    `[Go API Error] ${err.config?.method?.toUpperCase()} ${err.config?.url} [${err.response?.status}]:`,
                    err.response?.data
                );

                // ใช้ฟังก์ชันแปลงร่าง Error ของคุณ (เช่น แปลงเป็น ApiError หรือ ValidationError)
                const mapAxiosError = MapAxiosError(err);

                // 🚨 หมายเหตุ: ห้ามใส่ window.location.href ตรงนี้เด็ดขาด เพราะโค้ดรันบน Server 
                // หน้าที่การสั่งย้ายหน้าเมื่อ Unauthenticated (401) ควรส่งกลับไปให้ชั้น Action/Hook จัดการแทน

                throw mapAxiosError;
            }
        );
    }

    async get<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.get<T>(config.url, {
            params: config.params,
            headers: config.headers
        });
        return res.data;
    }

    async post<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.post<T>(
            config.url,
            config.body,
            {
                params: config.params,
                headers: config.headers
            }
        );
        return res.data;
    }

    async put<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.put<T>(
            config.url,
            config.body,
            {
                params: config.params,
                headers: config.headers
            }
        );
        return res.data;
    }

    async delete<T>(config: RequestConfig): Promise<T> {
        const res = await this.client.delete<T>(config.url, {
            params: config.params,
            headers: config.headers
        });
        return res.data;
    }
}