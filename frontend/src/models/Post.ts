import type {
    InternalFeaturesPostTransportHttpPaginationDTO as PaginationDTO,
    InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse as PostDTO,
    InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse as PostsPageDTO,
} from '@generated/data-contracts';

export interface PostAuthor {
    id: string;
    username: string;
    displayName: string;
    avatarUrl: string;
}

export interface Post {
    id: string;
    author: PostAuthor;
    parentId: string | null;
    content: string;
    imageUrl: string;
    likesCount: number;
    isLiked: boolean;
    createdAt: string;
    deletedAt: string | null;
}

export interface Pagination {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
}

export interface PostsPage {
    posts: Post[];
    pagination: Pagination;
}

export const toPost = (dto: PostDTO): Post => ({
    id: dto.id ?? '',
    author: {
        id: dto.user_id ?? '',
        username: dto.username ?? '',
        displayName: dto.display_name ?? '',
        avatarUrl: dto.avatar_url ?? '',
    },
    parentId: dto.parent_id ?? null,
    content: dto.content ?? '',
    imageUrl: dto.image_url ?? '',
    likesCount: dto.likes_count ?? 0,
    isLiked: dto.is_liked ?? false,
    createdAt: dto.created_at ?? '',
    deletedAt: dto.deleted_at ?? null,
});

export const toPagination = (dto: PaginationDTO | undefined): Pagination => ({
    page: dto?.page ?? 1,
    pageSize: dto?.page_size ?? 0,
    total: dto?.total ?? 0,
    totalPages: dto?.total_pages ?? 0,
});

export const toPostsPage = (dto: PostsPageDTO): PostsPage => ({
    posts: (dto.posts ?? []).map(toPost),
    pagination: toPagination(dto.pagination),
});
