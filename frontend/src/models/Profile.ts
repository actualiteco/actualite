import type {
    InternalFeaturesProfileTransportHttpProfileResponse,
    InternalFeaturesProfileTransportHttpUpdateProfileRequest,
} from '@generated/data-contracts';

type ProfileDTO = InternalFeaturesProfileTransportHttpProfileResponse & {
    is_following?: boolean | null;
};

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
    isFollowing: boolean | null;
}

export interface ProfileChanges {
    displayName?: string;
    username?: string;
    status?: string;
}

export const toProfile = (dto: ProfileDTO): Profile => ({
    id: dto.id ?? '',
    username: dto.username ?? '',
    displayName: dto.display_name ?? '',
    avatarUrl: dto.avatar_url ?? '',
    status: dto.status ?? '',
    email: '',
    followersCount: dto.followers_count ?? 0,
    followingCount: dto.following_count ?? 0,
    createdAt: '',
    isFollowing: dto.is_following ?? null,
});

export const toUpdateRequest = (changes: ProfileChanges) => {
    const request: Record<string, string> = {};

    if (changes.displayName !== undefined) {
        request.display_name = changes.displayName;
    }

    if (changes.username !== undefined) {
        request.username = changes.username;
    }

    if (changes.status !== undefined) {
        request.status = changes.status;
    }

    return request;
};

export const getDisplayName = (profile: Pick<Profile, 'displayName' | 'username'>): string =>
    profile.displayName.trim() || profile.username;
