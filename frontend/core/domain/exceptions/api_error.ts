// export type ApiErrorType =
//     | "NOT_FOUND"
//     | "VALIDATION_ERROR"
//     | "UNAUTHORIZED"
//     | "FORBIDDEN"
//     | "SERVER_ERROR"
//     | "NETWORK_ERROR"
//     | "UNKNOWN";

export class ApiError extends Error {
  constructor(
    public type: string,
    message: string,
    public statusCode?: number,
    public errorCode?: string,
    public metadata?: any
  ) {
    super(message);
  }
}