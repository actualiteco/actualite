import { create } from 'zustand';
import { ProfileService } from '@services/ProfileService';
import type { Profile } from '@models/Profile';
import { parseApiError } from '@utils/error.utils';

export type ProfileLoadStatus = 'idle' | 'loading' | 'success' | 'error';

export type ProfileLoadError = 'not_found' | 'forbidden' | 'network' | 'unknown';

interface ProfileState {
    profile: Profile | null;
    status: ProfileLoadStatus;
    error: ProfileLoadError | null;

    fetchProfile: (id?: string) => Promise<void>;

    setProfile: (profile: Profile) => void;

    reset: () => void;
}

type ProfileData = Pick<ProfileState, 'profile' | 'status' | 'error'>;

const initialState: ProfileData = {
    profile: null,
    status: 'idle',
    error: null,
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

export const useProfileStore = create<ProfileState>()((set) => ({
    ...initialState,

    async fetchProfile(id?: string) {
        const requestId = ++lastRequestId;
        set({ profile: null, status: 'loading', error: null });

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

    reset() {
        lastRequestId++;
        set({ ...initialState });
    },
}));
