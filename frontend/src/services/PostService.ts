import $api from '@api/api';
import type { InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse as PostsPageDTO } from '@generated/data-contracts';
import type { PostsPage } from '@models/Post';
import { toPostsPage } from '@models/Post';

export const POSTS_PAGE_SIZE = 20;

export class PostService {
    static async getByUser(
        userId: string,
        page = 1,
        pageSize = POSTS_PAGE_SIZE,
    ): Promise<PostsPage> {
        const response = await $api.get<PostsPageDTO>(
            `/profile/${encodeURIComponent(userId)}/posts`,
            { params: { page, page_size: pageSize } },
        );
        return toPostsPage(response.data);
    }
}
