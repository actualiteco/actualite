import { create } from 'zustand';
import { PostService } from '@services/PostService';
import type { Post } from '@models/Post';
import { parseApiError } from '@utils/error.utils';

export type UserPostsStatus = 'idle' | 'loading' | 'success' | 'error';
export type UserPostsError = 'network' | 'unknown';

interface UserPostsState {
    userId: string | null;
    posts: Post[];
    page: number;
    hasMore: boolean;
    status: UserPostsStatus;
    error: UserPostsError | null;
    isLoadingMore: boolean;
    loadMoreFailed: boolean;

    load: (userId: string) => Promise<void>;
    loadMore: () => Promise<void>;
    reset: () => void;
}

type UserPostsData = Omit<UserPostsState, 'load' | 'loadMore' | 'reset'>;

const initialState: UserPostsData = {
    userId: null,
    posts: [],
    page: 0,
    hasMore: false,
    status: 'idle',
    error: null,
    isLoadingMore: false,
    loadMoreFailed: false,
};

let lastRequestId = 0;

const toLoadError = (error: unknown): UserPostsError =>
    parseApiError(error).isNetwork ? 'network' : 'unknown';

const mergeUnique = (current: Post[], incoming: Post[]): Post[] => {
    const seen = new Set(current.map((post) => post.id));
    return [...current, ...incoming.filter((post) => !seen.has(post.id))];
};

export const useUserPostsStore = create<UserPostsState>()((set, get) => ({
    ...initialState,

    async load(userId: string) {
        const requestId = ++lastRequestId;
        set({ ...initialState, userId, status: 'loading' });

        try {
            const { posts, pagination } = await PostService.getByUser(userId, 1);
            if (requestId !== lastRequestId) {
                return;
            }
            set({
                posts: mergeUnique([], posts),
                page: 1,
                hasMore: posts.length > 0 && pagination.page < pagination.totalPages,
                status: 'success',
            });
        } catch (error) {
            if (requestId !== lastRequestId) {
                return;
            }
            set({ status: 'error', error: toLoadError(error) });
        }
    },

    async loadMore() {
        const { userId, status, hasMore, isLoadingMore, page } = get();
        if (!userId || status !== 'success' || !hasMore || isLoadingMore) {
            return;
        }

        const requestId = ++lastRequestId;
        const nextPage = page + 1;
        set({ isLoadingMore: true, loadMoreFailed: false });

        try {
            const { posts, pagination } = await PostService.getByUser(userId, nextPage);
            if (requestId !== lastRequestId) {
                return;
            }
            set((state) => ({
                posts: mergeUnique(state.posts, posts),
                page: nextPage,
                hasMore: posts.length > 0 && pagination.page < pagination.totalPages,
                isLoadingMore: false,
            }));
        } catch {
            if (requestId !== lastRequestId) {
                return;
            }
            set({ isLoadingMore: false, loadMoreFailed: true });
        }
    },

    reset() {
        lastRequestId++;
        set({ ...initialState });
    },
}));
