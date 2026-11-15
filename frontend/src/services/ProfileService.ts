import $api from '@api/api';
import type { InternalFeaturesProfileTransportHttpProfileResponse as ProfileDTO } from '@generated/data-contracts';
import type { Profile } from '@models/Profile';
import { toProfile } from '@models/Profile';

export class ProfileService {
    static async getMe(): Promise<Profile> {
        const response = await $api.get<ProfileDTO>('/profile/me');
        return toProfile(response.data);
    }

    static async getById(id: string): Promise<Profile> {
        const response = await $api.get<ProfileDTO>(`/profile/${encodeURIComponent(id)}`);
        return toProfile(response.data);
    }
}
