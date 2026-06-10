export interface CookieSessionStoragePort {
    save(key: string, value: string, maxAgeInSeconds: number): Promise<void>;
    get(key: string): Promise<string | undefined>;
    clearSession(key: string): Promise<void>;
}