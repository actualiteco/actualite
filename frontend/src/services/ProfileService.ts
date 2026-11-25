import $api from '@api/api';
import { toProfile, toUpdateRequest, type Profile, type ProfileChanges } from '@models/Profile';
import type {
    InternalFeaturesProfileTransportHttpProfileResponse,
    InternalFeaturesProfileTransportHttpUploadAvatarResponse,
} from '@generated/data-contracts';

export class ProfileService {
    static async getMe(): Promise<Profile> {
        const { data } = await $api.get<InternalFeaturesProfileTransportHttpProfileResponse>('/profile/me');
        return toProfile(data);
    }

    static async getById(userId: string): Promise<Profile> {
        const { data } = await $api.get<InternalFeaturesProfileTransportHttpProfileResponse>(`/profile/${userId}`);
        return toProfile(data);
    }

    static async update(changes: ProfileChanges): Promise<Profile> {
        const { data } = await $api.patch<InternalFeaturesProfileTransportHttpProfileResponse>(
            '/profile',
            toUpdateRequest(changes),
        );
        return toProfile(data);
    }

    static async uploadAvatar(file: Blob, fileName: string): Promise<string> {
        const form = new FormData();
        form.append('file', file, fileName);
        const { data } =
            await $api.post<InternalFeaturesProfileTransportHttpUploadAvatarResponse>(
                '/profile/avatar',
                form,
            );
        if (!data.avatar_url) {
            throw new Error('Сервер не вернул адрес аватара');
        }
        return data.avatar_url;
    }

    static async follow(userId: string): Promise<void> {
        await $api.post(`/profile/${userId}/follow`);
    }

    static async unfollow(userId: string): Promise<void> {
        await $api.post(`/profile/${userId}/unfollow`);
    }

    static async isUsernameAvailable(username: string): Promise<boolean> {
        try {
            await $api.get('/profile/username-availability', { params: { username } });
            return true;
        } catch (error) {
            if ((error as { response?: { status?: number } }).response?.status === 409) {
                return false;
            }
            throw error;
        }
    }
}
