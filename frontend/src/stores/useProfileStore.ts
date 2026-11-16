import { create } from 'zustand';
import { ProfileService } from '@services/ProfileService';
import type { Profile, ProfileChanges } from '@models/Profile';
import { useAuthStore } from '@stores/useAuthStore';
import { parseApiError } from '@utils/error.utils';

export type ProfileLoadStatus = 'idle' | 'loading' | 'success' | 'error';

export type ProfileLoadError = 'not_found' | 'forbidden' | 'network' | 'unknown';

interface ProfileState {
    profile: Profile | null;
    status: ProfileLoadStatus;
    error: ProfileLoadError | null;
    isFollowPending: boolean;

    fetchProfile: (id?: string) => Promise<void>;
    setProfile: (profile: Profile) => void;
    updateProfile: (changes: ProfileChanges) => Promise<void>;
    toggleFollow: () => Promise<void>;
    reset: () => void;
}

type ProfileData = Pick<ProfileState, 'profile' | 'status' | 'error' | 'isFollowPending'>;

const initialState: ProfileData = {
    profile: null,
    status: 'idle',
    error: null,
    isFollowPending: false,
};

let lastRequestId = 0;

const toLoadError = (error: unknown): ProfileLoadError => {
    const { status, isNetwork } = parseApiError(error);

    if (isNetwork) {
        return 'network';
    }
    if (status === 404 || status === 400) {
        return 'not_found';
    }
    if (status === 403) {
        return 'forbidden';
    }
    return 'unknown';
};

export const useProfileStore = create<ProfileState>()((set, get) => ({
    ...initialState,

    async fetchProfile(id?: string) {
        const requestId = ++lastRequestId;
        set({ profile: null, status: 'loading', error: null, isFollowPending: false });

        try {
            const profile = id ? await ProfileService.getById(id) : await ProfileService.getMe();
            if (requestId !== lastRequestId) {
                return;
            }
            set({ profile, status: 'success', error: null });
        } catch (error) {
            if (requestId !== lastRequestId) {
                return;
            }
            set({ profile: null, status: 'error', error: toLoadError(error) });
        }
    },

    setProfile(profile: Profile) {
        set({ profile, status: 'success', error: null });
    },

    async updateProfile(changes: ProfileChanges) {
        const updated = await ProfileService.update(changes);
        set({ profile: updated, status: 'success', error: null });
        useAuthStore.setState({ username: updated.username });
    },

    async toggleFollow() {
        const { profile, isFollowPending } = get();
        if (!profile || isFollowPending) {
            return;
        }

        const targetId = profile.id;
        const wasFollowing = profile.isFollowing === true;
        set({ isFollowPending: true });

        const finish = (isFollowing: boolean | null, delta: number) =>
            set((state) => {
                if (!state.profile || state.profile.id !== targetId) {
                    return {};
                }
                if (isFollowing === null) {
                    return { isFollowPending: false };
                }
                return {
                    isFollowPending: false,
                    profile: {
                        ...state.profile,
                        isFollowing,
                        followersCount: Math.max(0, state.profile.followersCount + delta),
                    },
                };
            });

        try {
            if (wasFollowing) {
                await ProfileService.unfollow(targetId);
            } else {
                await ProfileService.follow(targetId);
            }
            finish(!wasFollowing, wasFollowing ? -1 : 1);
        } catch (error) {
            if (parseApiError(error).status === 400) {
                finish(!wasFollowing, 0);
                return;
            }
            finish(null, 0);
            throw error;
        }
    },

    reset() {
        lastRequestId++;
        set({ ...initialState });
    },
}));
