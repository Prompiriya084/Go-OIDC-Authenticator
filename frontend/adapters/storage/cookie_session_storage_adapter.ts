import { cookies } from 'next/headers';
import { CookieSessionStoragePort } from '@/core/ports/storage/cookie_session_storage_port';

export class CookieSessionAdapter implements CookieSessionStoragePort {
  async saveToken(token: string, maxAgeInSeconds: number): Promise<void> {
    const cookieStore = await cookies();
    cookieStore.set('auth_session', token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: maxAgeInSeconds,
    });
  }

  async getToken(): Promise<string | undefined> {
    const cookieStore = await cookies();
    return cookieStore.get('auth_session')?.value;
  }

  async clearSession(): Promise<void> {
    const cookieStore = await cookies();
    cookieStore.delete('auth_session');
  }
}