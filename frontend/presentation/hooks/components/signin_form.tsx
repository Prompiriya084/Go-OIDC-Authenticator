import { SigninRequestDTO } from "@/core/dtos/signin_request_dto";
import { useSearchParams } from "next/navigation";
import { useAccount } from "../use_account";
import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import { FieldError } from "@/components/ui/field";
export function SigninForm({
    className,
    ...props
}: React.ComponentProps<"div">) {
    const searchParams = useSearchParams();
    const { submitSignIn, loading, fieldError } = useAccount();

    async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const getParmas = (key: string): string => searchParams.get(key) ?? "";
        const authParams: AuthQueryParameterDTO = {
            flowId: getParmas('flowId'),
            clientId: getParmas("clientId")
        };

        const formData = new FormData(e.currentTarget);
        const signinInput: SigninRequestDTO = {
            username: String(formData.get("username")),
            password: String(formData.get("password"))
        }
        await submitSignIn(authParams, signinInput);
    }

    return (
        <div className={cn("flex flex-col gap-6", className)} {...props}>
            <div className="w-full max-w-4xl flex flex-col md:flex-row bg-white rounded-xl shadow-[0px_8px_32px_rgba(0,0,0,0.06)] overflow-hidden border border-gray-100 max-h-[85vh] md:h-auto">
                {/* <!-- Left Zone: How It Works Area --> */}
                <section className="hidden md:flex md:w-[42%] bg-surface-container items-center justify-center p-8 lg:p-10 border-r border-gray-100 overflow-hidden">
                    <div className="w-full max-w-[280px]">
                        <h2 className="text-[11px] font-label-bold text-primary uppercase tracking-widest mb-6">How it works</h2>
                        {/* <!-- Animated Illustration Area --> */}
                        <div className="relative h-48 mb-8 bg-surface-container-low rounded-xl border border-dashed border-outline-variant flex items-center justify-center overflow-hidden">
                            <div className="step-container w-full h-full">
                                {/* <!-- Step 1 Illustration: Credentials --> */}
                                <div className="animate-step-1">
                                    <div className="w-40 bg-white p-4 rounded-lg shadow-sm border border-gray-100 space-y-3">
                                        <div className="space-y-1">
                                            <div className="h-1 w-12 bg-gray-200 rounded"></div>
                                            <div className="h-6 w-full bg-gray-50 border border-gray-100 rounded overflow-hidden relative flex items-center px-2">
                                                <div className="text-[9px] font-medium text-primary animate-typing-user overflow-hidden whitespace-nowrap">example</div>
                                            </div>
                                        </div>
                                        <div className="space-y-1">
                                            <div className="h-1 w-12 bg-gray-200 rounded"></div>
                                            <div className="h-6 w-full bg-gray-50 border border-gray-100 rounded overflow-hidden relative flex items-center px-2">
                                                <div className="flex gap-1 animate-typing-pass opacity-0">
                                                    <div className="w-1 h-1 rounded-full bg-gray-400"></div>
                                                    <div className="w-1 h-1 rounded-full bg-gray-400"></div>
                                                    <div className="w-1 h-1 rounded-full bg-gray-400"></div>
                                                    <div className="w-1 h-1 rounded-full bg-gray-400"></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div className="h-6 w-full bg-primary rounded flex items-center justify-center">
                                            <div className="w-10 h-1 bg-white/20 rounded-full"></div>
                                        </div>
                                    </div>
                                </div>
                                {/* <!-- Step 2 Illustration: Line Art Phone with Sliding Code --> */}
                                <div className="animate-step-2">
                                    <div className="relative w-24 h-40 bg-white rounded-xl border-2 border-black flex flex-col p-1.5 overflow-hidden shadow-sm">
                                        <div className="w-6 h-0.5 bg-black rounded-full mx-auto mb-4"></div>
                                        <div className="flex-grow relative flex flex-col items-center justify-center">
                                            <div className="animate-phone-code absolute bottom-0 flex flex-col items-center w-full">
                                                <div className="w-full px-1.5">
                                                    <div className="bg-[#2e3132] rounded-md py-1.5 flex flex-col items-center justify-center text-white shadow-md">
                                                        <div className="text-[12px] font-bold tracking-widest">123 456</div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div className="w-8 h-0.5 bg-black rounded-full mx-auto mt-auto mb-0.5"></div>
                                    </div>
                                </div>
                                {/* <!-- Step 3 Illustration: Notebook Access Granted --> */}
                                <div className="animate-step-3">
                                    <div className="laptop-wrapper flex flex-col items-center w-full scale-[0.85]">
                                        <div className="w-56 h-36 bg-gray-200 rounded-t-xl border-t-2 border-x-2 border-gray-300 relative flex items-center justify-center overflow-hidden">
                                            <div className="animate-screen-on absolute inset-0 bg-white flex flex-col items-center justify-center">
                                                <div className="animate-verified-check-new flex flex-col items-center justify-center">
                                                    <div className="bg-white p-4 rounded-xl shadow-[0_8px_20px_rgba(0,0,0,0.1)] border border-gray-100 flex flex-col items-center gap-3">
                                                        <div className="w-12 h-12 rounded-full bg-green-50 flex items-center justify-center ring-4 ring-green-100">
                                                            <span className="material-symbols-outlined text-green-500 font-bold text-3xl" data-icon="check_circle">check_circle</span>
                                                        </div>
                                                        <div className="text-center">
                                                            <div className="text-[11px] font-black text-gray-900 uppercase tracking-widest">Access Granted</div>
                                                            <div className="text-[9px] text-gray-400 mt-0.5">Workspace Ready</div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="absolute top-1.5 left-1/2 -translate-x-1/2 w-0.5 h-0.5 bg-gray-400 rounded-full"></div>
                                        </div>
                                        <div className="w-64 h-2.5 bg-gray-400 rounded-b-xl shadow-lg relative">
                                            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-12 h-0.5 bg-gray-500/50 rounded-b-md"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        {/* <!-- Step Indicators --> */}
                        <div className="space-y-5">
                            <div className="flex items-start gap-3">
                                <div className="flex-shrink-0 w-7 h-7 bg-surface-container relative flex items-center justify-center rounded-full">
                                    <span className="material-symbols-outlined text-[13px] text-gray-400" data-icon="login">login</span>
                                    <div className="absolute inset-0 bg-primary-fixed flex items-center justify-center animate-step-1 opacity-0">
                                        <span className="material-symbols-outlined text-[13px] text-primary" data-icon="login">login</span>
                                    </div>
                                    <div className="absolute top-full left-1/2 -translate-x-1/2 w-[1px] h-5 bg-gray-200">
                                        <div className="bg-primary w-full h-0 animate-line-1"></div>
                                    </div>
                                </div>
                                <div>
                                    <h3 className="text-xs font-label-bold text-on-surface">Step 1: Credentials</h3>
                                    <p className="text-[10px] text-secondary mt-0.5">Enter your Username and Password.</p>
                                </div>
                            </div>
                            <div className="flex items-start gap-3">
                                <div className="flex-shrink-0 w-7 h-7 bg-surface-container relative flex items-center justify-center rounded-full">
                                    <span className="material-symbols-outlined text-[13px] text-gray-400" data-icon="smartphone">smartphone</span>
                                    <div className="absolute inset-0 bg-primary-fixed flex items-center justify-center animate-step-2 opacity-0">
                                        <span className="material-symbols-outlined text-[13px] text-primary" data-icon="smartphone">smartphone</span>
                                    </div>
                                    <div className="absolute top-full left-1/2 -translate-x-1/2 w-[1px] h-5 bg-gray-200">
                                        <div className="bg-primary w-full h-0 animate-line-2"></div>
                                    </div>
                                </div>
                                <div>
                                    <h3 className="text-xs font-label-bold text-on-surface">Step 2: Verification</h3>
                                    <p className="text-[10px] text-secondary mt-0.5">Open the authenticator app on your device for security code.</p>
                                </div>
                            </div>
                            <div className="flex items-start gap-3">
                                <div className="flex-shrink-0 w-7 h-7 bg-surface-container relative overflow-hidden flex items-center justify-center rounded-full">
                                    <span className="material-symbols-outlined text-[13px] text-gray-400" data-icon="done_all">done_all</span>
                                    <div className="absolute inset-0 bg-primary-fixed flex items-center justify-center animate-step-3 opacity-0">
                                        <span className="material-symbols-outlined text-[13px] text-primary" data-icon="done_all">done_all</span>
                                    </div>
                                </div>
                                <div>
                                    <h3 className="text-xs font-label-bold text-on-surface">Step 3: Access</h3>
                                    <p className="text-[10px] text-secondary mt-0.5">Enter code to complete login and access your workspace.</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
                {/* <!-- Right Zone: Interaction/Login Form Area --> */}
                <section className="flex-grow md:w-[58%] flex items-center justify-center p-8 lg:p-12 bg-white overflow-y-auto custom-scrollbar">
                    <div className="w-full max-w-[320px] flex flex-col gap-5">
                        <div>
                            <h1 className="font-headline-md text-headline-md text-on-surface mb-1.5">Welcome Back</h1>
                            <p className="font-body-sm text-body-sm text-on-surface-variant">Enter your credentials for secure enterprise access.</p>
                        </div>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="space-y-1.5">
                                <label className="font-label-bold text-on-surface">Username {fieldError.username && <FieldError><p className="text-red-500">*{fieldError.username}</p></FieldError>}</label>
                                <div className="relative">
                                    <Input
                                        className="w-full h-[44px] px-3.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all placeholder:text-gray-400 text-sm"
                                        placeholder="e.g. j.doe@enterprise.com"
                                        type="text"
                                        name="username"
                                        disabled={loading} />
                                    <span className="material-symbols-outlined absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-400 text-xl" data-icon="person">person</span>
                                </div>
                            </div>
                            <div className="space-y-1.5">
                                <div className="flex justify-between items-center">
                                    <label className="font-label-bold text-on-surface" >Password {fieldError.password && <span className="text-red-500 text-sm">*{fieldError.password}</span>}</label>
                                    {/* <a className="text-[11px] text-gray-500 hover:text-black underline underline-offset-2 transition-colors" href="#">Forgot password?</a> */}
                                </div>
                                <div className="relative">
                                    <Input
                                        className="w-full h-[44px] px-3.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent outline-none transition-all placeholder:text-gray-400 text-sm"
                                        name="password"
                                        placeholder="••••••••••••"
                                        type="password"
                                        disabled={loading} />
                                    <span className="material-symbols-outlined absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-400 text-xl" data-icon="lock">lock</span>
                                </div>
                            </div>
                            {/* <div className="flex items-center gap-2.5 py-1">
                <Input className="w-3.5 h-3.5 text-primary border-gray-300 rounded focus:ring-primary" id="remember" type="checkbox" />
                <label className="text-[12px] text-on-surface">Keep me logged in for 8 hours</label>
              </div> */}
                            <Button
                                className="w-full h-[44px] bg-primary text-on-primary font-button-text rounded-lg hover:opacity-90 active:scale-[0.98] transition-all flex items-center justify-center gap-2"
                                type="submit"
                                disabled={loading}>
                                {loading && <Spinner data-icon="inline-startt" />}
                                {loading ? "loading...." : "Continue to 2FA"}
                                {!loading && <span className="material-symbols-outlined text-lg" data-icon="arrow_forward">arrow_forward</span>}
                            </Button>

                            <div className="flex flex-col gap-2 pt-4 border-t border-gray-100">
                                <button className="flex items-center justify-center gap-2 text-on-surface-variant hover:text-primary transition-colors font-body-sm text-[13px]" type="button">
                                    <span className="material-symbols-outlined text-sm" data-icon="help_outline">help_outline</span>
                                    {/* Trouble scanning? */}
                                    <a className="text-gray-500 hover:text-black underline underline-offset-2 transition-colors" href="#">Forgot password?</a>
                                </button>
                            </div>
                        </form>
                    </div>
                </section>
            </div>
        </div>
    )
}