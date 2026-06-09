export interface CookieSessionStoragePort {
    saveToken(token: string, maxAgeInSeconds: number): Promise<void>;
    getToken(): Promise<string | undefined>;
    clearSession(): Promise<void>;
}