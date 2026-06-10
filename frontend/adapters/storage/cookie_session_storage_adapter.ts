import { cookies } from 'next/headers';
import { CookieSessionStoragePort } from '@/core/ports/storage/cookie_session_storage_port';

export class CookieSessionAdapter implements CookieSessionStoragePort {
  
    async save(key: string, value: string, maxAgeInSeconds: number): Promise<void> {
    const cookieStore = await cookies();
    cookieStore.set(key, value, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: maxAgeInSeconds,
    });
  }

  async get(key: string): Promise<string | undefined> {
    const cookieStore = await cookies();
    const cookie = cookieStore.get(key);
    return cookie ? cookie.value : undefined;
  }

  async clearSession(key: string): Promise<void> {
    const cookieStore = await cookies();
    cookieStore.delete(key);
  }

}