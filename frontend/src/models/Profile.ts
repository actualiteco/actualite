import type { InternalFeaturesProfileTransportHttpProfileResponse as ProfileDTO } from '@generated/data-contracts';

export interface Profile {
    id: string;
    username: string;
    displayName: string;
    avatarUrl: string;
    status: string;
    email: string;
    followersCount: number;
    followingCount: number;
    createdAt: string;
}

export const toProfile = (dto: ProfileDTO): Profile => ({
    id: dto.id ?? '',
    username: dto.username ?? '',
    displayName: dto.display_name ?? '',
    avatarUrl: dto.avatar_url ?? '',
    status: dto.status ?? '',
    email: dto.email ?? '',
    followersCount: dto.followers_count ?? 0,
    followingCount: dto.following_count ?? 0,
    createdAt: dto.created_at ?? '',
});

export const getDisplayName = (profile: Pick<Profile, 'displayName' | 'username'>): string =>
    profile.displayName.trim() || profile.username;
