import $api from '@api/api';
import type {
    InternalFeaturesProfileTransportHttpProfileResponse as ProfileDTO,
    InternalFeaturesProfileTransportHttpUploadAvatarResponse as UploadAvatarDTO,
} from '@generated/data-contracts';
import type { Profile, ProfileChanges } from '@models/Profile';
import { toProfile, toUpdateRequest } from '@models/Profile';

export class ProfileService {
    static async getMe(): Promise<Profile> {
        const response = await $api.get<ProfileDTO>('/profile/me');
        return toProfile(response.data);
    }

    static async getById(id: string): Promise<Profile> {
        const response = await $api.get<ProfileDTO>(`/profile/${encodeURIComponent(id)}`);
        return toProfile(response.data);
    }

    static async update(changes: ProfileChanges): Promise<Profile> {
        const response = await $api.patch<ProfileDTO>('/profile', toUpdateRequest(changes));
        return toProfile(response.data);
    }

    static async uploadAvatar(file: Blob): Promise<string> {
        const formData = new FormData();
        formData.append('file', file, file.type === 'image/jpeg' ? 'avatar.jpg' : 'avatar.png');
        const response = await $api.post<UploadAvatarDTO>('/profile/avatar', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        });
        return response.data.avatar_url ?? '';
    }

    static async follow(id: string): Promise<void> {
        await $api.post(`/profile/${encodeURIComponent(id)}/follow`);
    }

    static async unfollow(id: string): Promise<void> {
        await $api.post(`/profile/${encodeURIComponent(id)}/unfollow`);
    }
}
