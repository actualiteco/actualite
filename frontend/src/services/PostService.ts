import $api from '@api/api';
import type { CreatePostRequest, CreatedPost, Post, PostsPage } from '@models/api.types';
import type { InternalFeaturesPostTransportHttpImageUrlDTOResponse } from '@generated/data-contracts';

const pageParams = (page: number, pageSize?: number): Record<string, number> => {
  const params: Record<string, number> = { page };
  if (pageSize !== undefined && pageSize !== null && pageSize > 0) {
    params.page_size = pageSize;
  }
  return params;
};

export class PostService {
    static async getFeed(page: number, pageSize: number): Promise<PostsPage> {
        const { data } = await $api.get<PostsPage>('/posts', { params: pageParams(page, pageSize) });
        return data;
    }

    static async getById(postId: string): Promise<Post> {
        const { data } = await $api.get<Post>(`/posts/${postId}`);
        return data;
    }

    static async getReplies(postId: string, page: number, pageSize: number): Promise<PostsPage> {
        const { data } = await $api.get<PostsPage>(`/posts/${postId}/replies`, {
            params: pageParams(page, pageSize),
        });
        return data;
    }

    static async getUserPosts(userId: string, page: number, pageSize?: number): Promise<PostsPage> {
        const { data } = await $api.get<PostsPage>(`/profile/${userId}/posts`, {
            params: pageParams(page, pageSize),
        });
        return data;
    }

    static async create(payload: CreatePostRequest): Promise<CreatedPost> {
        const { data } = await $api.post<CreatedPost>('/posts/create', payload);
        return data;
    }

    static async uploadImage(file: File): Promise<string> {
        const form = new FormData();
        form.append('file', file);
        const { data } = await $api.post<InternalFeaturesPostTransportHttpImageUrlDTOResponse>(
            '/posts/image',
            form,
        );
        if (!data.image_url) {
            throw new Error('Сервер не вернул адрес изображения');
        }
        return data.image_url;
    }

    static async remove(postId: string): Promise<void> {
        await $api.delete(`/posts/${postId}`);
    }

    static async like(postId: string): Promise<void> {
        await $api.post(`/posts/${postId}/like`);
    }

    static async unlike(postId: string): Promise<void> {
        await $api.post(`/posts/${postId}/unlike`);
    }
}
