import $api from '@api/api';
import type {
    LoginRequest,
    RegisterRequest,
    RegisterResponse,
    TokensResponse,
} from '@models/api.types';

export class AuthService {
    static async register(data: RegisterRequest): Promise<RegisterResponse> {
        const response = await $api.post<RegisterResponse>('auth/register', data);
        return response.data;
    }

    static async login(data: LoginRequest): Promise<TokensResponse> {
        const response = await $api.post<TokensResponse>('auth/login', data);
        return response.data;
    }

    static async logout(accessToken: string): Promise<void> {
        await $api.post('/auth/logout', null, {
            headers: { Authorization: `Bearer ${accessToken}` },
        });
    }
}
