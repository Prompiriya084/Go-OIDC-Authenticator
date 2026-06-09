export interface HttpPort {
    get<T>(config: RequestConfig): Promise<T>;
    post<T>(config: RequestConfig): Promise<T>;
    put<T>(config: RequestConfig): Promise<T>;
    delete<T>(config: RequestConfig): Promise<T>;
}
export interface RequestConfig {
//   method: "GET" | "POST" | "PUT" | "DELETE";
  url: string;
  body?: any;
  params?: any;
  // tokenType: TokenType;
}