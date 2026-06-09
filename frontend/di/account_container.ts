import { AxiosHttpAdapter } from '@/adapters/http/axios_http_adapter';
import { AccountApiRepositoryAdapter } from '@/adapters/repositories/account_api_repository_adapter';
import { CookieSessionAdapter } from '@/adapters/storage/cookie_session_storage_adapter';
import { AccountService } from '@/core/services/account_service';

/**
 * DI Container
 */
class AccountDIContainer {
    private static instance: AccountDIContainer;

    private readonly axiosHttp: AxiosHttpAdapter
    private readonly accountApiRepository: AccountApiRepositoryAdapter;
    public readonly accountService: AccountService;
    public readonly cookieSessionStorage: CookieSessionAdapter;

    private constructor() {
        this.axiosHttp = new AxiosHttpAdapter("localhost:8080")
        this.accountApiRepository = new AccountApiRepositoryAdapter(this.axiosHttp);
        this.accountService = new AccountService(this.accountApiRepository);

        this.cookieSessionStorage = new CookieSessionAdapter();
    }

    public static getInstance(): AccountDIContainer {
        if (!AccountDIContainer.instance) {
            AccountDIContainer.instance = new AccountDIContainer();
        }
        return AccountDIContainer.instance;
    }
}

// Export ตัวแปรพร้อมใช้ไปให้เลเยอร์อื่นเรียกใช้งานได้ทันที
export const accountContainer = AccountDIContainer.getInstance();