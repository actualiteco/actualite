import $api from '@api/api';
import type { InternalFeaturesProfileTransportHttpProfileResponse as ProfileDTO } from '@generated/data-contracts';
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

    static async follow(id: string): Promise<void> {
        await $api.post(`/profile/${encodeURIComponent(id)}/follow`);
    }

    static async unfollow(id: string): Promise<void> {
        await $api.post(`/profile/${encodeURIComponent(id)}/unfollow`);
    }
}
