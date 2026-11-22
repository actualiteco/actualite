/** Короткие имена типов API. Меняются только здесь, если генератор переименует схемы. */
import type {
    InternalFeaturesAuthTransportHttpLoginDTO,
    InternalFeaturesAuthTransportHttpRegisterDTO,
    InternalFeaturesAuthTransportHttpRegisterDTOResponse,
    InternalFeaturesAuthTransportHttpTokensDTOResponse,
    InternalFeaturesPostTransportHttpCreatePostDTO,
    InternalFeaturesPostTransportHttpPostDTOResponse,
    InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse,
    InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
    InternalFeaturesProfileTransportHttpProfileResponse,
    InternalFeaturesProfileTransportHttpUpdateProfileRequest,
} from '@generated/data-contracts';

export type LoginRequest = InternalFeaturesAuthTransportHttpLoginDTO;
export type RegisterRequest = InternalFeaturesAuthTransportHttpRegisterDTO;
export type RegisterResponse = InternalFeaturesAuthTransportHttpRegisterDTOResponse;
export type TokensResponse = InternalFeaturesAuthTransportHttpTokensDTOResponse;

export type CreatePostRequest = InternalFeaturesPostTransportHttpCreatePostDTO;
export type CreatedPost = InternalFeaturesPostTransportHttpPostDTOResponse;

export type Post = InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse;
export type PostsPage = InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse;

export type Profile = InternalFeaturesProfileTransportHttpProfileResponse;
export type UpdateProfileRequest = InternalFeaturesProfileTransportHttpUpdateProfileRequest;
