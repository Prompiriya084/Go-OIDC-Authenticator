
import { ApiError } from "@/core/domain/exceptions/api_error";
import axios from "axios";

export function MapAxiosError(error: unknown): ApiError {
    if (!axios.isAxiosError(error)) {
        return new ApiError("UNKNOWN", "Unexpected error");
    }

    if (!error.response) {
        if ((error as any).code === "ECONNABORTED") {
            return new ApiError("NETWORK_ERROR", "Request timeout", 408);
        }

        return new ApiError("NETWORK_ERROR", "Network error");
    }

    const status = error.response.status;
    const data = error.response.data as any;
    const errorCode = data?.error;

    const metadata = data;

    // --- START ERROR EXTRACTION ---
    let message = "Request failed";

    if (data?.errors) {
        // 1. Get the keys (e.g., ["clientId", "flowId"])
        const errorKeys = Object.keys(data.errors);

        if (errorKeys.length > 0) {
            // 2. Get the first field's array (e.g., ["The clientId field is required."])
            const firstFieldErrors = data.errors[errorKeys[0]];

            // 3. Get the first string in that array
            message = firstFieldErrors[0];
        }
    } else {
        message = data?.title || data?.error_description || data?.message || error.message;
    }

    switch (status) {
        case 400:
            return new ApiError("VALIDATION_ERROR", message, status, errorCode, metadata);
        case 401:
            return new ApiError("UNAUTHORIZED", message, status, errorCode, metadata);
        case 403:
            return new ApiError("FORBIDDEN", message, status, errorCode, metadata);
        case 404:
            return new ApiError("NOT_FOUND", message, status, errorCode, metadata);
        case 410:
            return new ApiError("SESSION_EXPIRED", message, status, errorCode, metadata);
        case 500:
            return new ApiError("SERVER_ERROR", message, status, errorCode, metadata);
        default:
            return new ApiError("UNKNOWN", message, status, errorCode, metadata);
    }
}