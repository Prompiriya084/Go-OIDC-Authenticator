import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { SessionKey } from './core/domain/constants/session_keys';

export function middleware(request: NextRequest) {
  const { pathname, searchParams } = request.nextUrl;

  const hasPreMfaToken = request.cookies.has(SessionKey.PreMfa); 
  const hasMfaToken = request.cookies.has(SessionKey.Mfa);

  const flowId = searchParams.get('flowId') || '';
  const clientId = searchParams.get('clientId') || '';

  const preserveOidcParams = (url: URL) => {
    url.searchParams.set('flowId', flowId);
    url.searchParams.set('clientId', clientId);
    return url;
  };
  // 🚨 ด่านที่ A: ป้องกันโซนเซ็ตอัพ MFA (First-time Register / Setup Totp)
  if (pathname.startsWith('/mfa/setup-totp') || pathname.startsWith('/mfa/confirm-totp')) {
    // ถ้าไม่มีตั๋ว pre_mfa_token ติดตัวมาเลย แต่อุตริจะเข้าหน้าเซ็ตอัพ
    if (!hasPreMfaToken) {
      const signInUrl = new URL('/account/signin', request.url);
      return NextResponse.redirect(preserveOidcParams(signInUrl));
    }
    return NextResponse.next();
  }

  // 🚨 ด่านที่ B: ป้องกันโซนยืนยันรหัส MFA ประจำวัน (Verify Totp)
  if (pathname.startsWith('/mfa/verify-totp')) {
    // ถ้าไม่มีตั๋ว mfa_stage_token ติดตัวมาเลย แต่อยากจะวิ่งเข้ามากรอกรหัส 6 ตัว
    if (!hasMfaToken) {
      const signInUrl = new URL('/account/signin', request.url);
      return NextResponse.redirect(preserveOidcParams(signInUrl));
    }
    return NextResponse.next();
  }

  return NextResponse.next();
}

export const config = {
  // ควบคุมทุกพาร์ทที่ขึ้นต้นด้วย /mfa เพื่อไม่ให้ Middleware ทำงานหนักเกินไปในหน้า Static Assets
  matcher: ['/mfa/:path*'],
};