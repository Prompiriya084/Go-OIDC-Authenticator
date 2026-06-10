import { AxiosHttpAdapter } from '@/adapters/http/axios_http_adapter';
import { AccountApiRepositoryAdapter } from '@/adapters/repositories/account_api_repository_adapter';
import { MfaApiRepositoryAdapter } from '@/adapters/repositories/mfa_api_repository_adapter';
import { CookieSessionAdapter } from '@/adapters/storage/cookie_session_storage_adapter';
import { AccountService } from '@/core/services/account_service';
import { MfaService } from '@/core/services/mfa_service';

/**
 * DI Container
 */
class MfaDIContainer {
    private static instance: MfaDIContainer;

    private readonly axiosHttp: AxiosHttpAdapter
    private readonly mfaApiRepository: MfaApiRepositoryAdapter;
    public readonly mfaService: MfaService;
    public readonly cookieSessionStorage: CookieSessionAdapter

    private constructor() {
        this.axiosHttp = new AxiosHttpAdapter("localhost:8080")
        this.mfaApiRepository = new MfaApiRepositoryAdapter(this.axiosHttp);
        this.mfaService = new MfaService(this.mfaApiRepository);
        this.cookieSessionStorage = new CookieSessionAdapter();
    }

    public static getInstance(): MfaDIContainer {
        if (!MfaDIContainer.instance) {
            MfaDIContainer.instance = new MfaDIContainer();
        }
        return MfaDIContainer.instance;
    }
}

// Export ตัวแปรพร้อมใช้ไปให้เลเยอร์อื่นเรียกใช้งานได้ทันที
export const mfaContainer = MfaDIContainer.getInstance();