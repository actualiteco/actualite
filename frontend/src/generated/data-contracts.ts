/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ActualitecoActualiteInternalCoreDomainCustomError {
  code?: string;
  /** @example "invalid request" */
  message?: string;
}

export interface ActualitecoActualiteInternalCoreDomainInternalError {
  code?: "INTERNAL_ERROR";
  message?: string;
}

export interface InternalFeaturesAuthTransportHttpLoginDTO {
  /** @example "test.user@email.com" */
  email: string;
  /**
   * @minLength 8
   * @maxLength 64
   * @example "hashpassword"
   */
  password: string;
}

export interface InternalFeaturesAuthTransportHttpLoginDTOResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTM2NjUxNzcsImlhdCI6MTc5MzYyMTk3NywiVXNlcklkIjoiODc5NGU1ODktMGRkYi00M2NlLTlmOTItMTZmYWFmY2Y0ZWU0IiwiUm9sZSI6InVzZXIifQ.iv_JT1Fke-QaK2nx2FEF_pkDroyJP-_xK81egOAV7h4" */
  token?: string;
  /** @example "8794e589-0ddb-43ce-9f92-16faafcf4ee4" */
  user_id?: string;
  /** @example "username" */
  username?: string;
}

export interface InternalFeaturesAuthTransportHttpRegisterDTO {
  /** @example "m4rkek@gmail.com" */
  email: string;
  /** @example "m4rkek" */
  username?: string;
  /**
   * @minLength 8
   * @maxLength 64
   * @example "supersecretpass"
   */
  password: string;
  /** @example "admin" */
  role: string;
}

export interface InternalFeaturesAuthTransportHttpRegisterDTOResponse {
  /** @example "3fa85f64-5717-4562-b3fc-2c963f66afa6" */
  id?: string;
  /** @example "m4rkek" */
  username?: string;
  /** @example "m4rkek@gmail.com" */
  email?: string;
  /** @example "user" */
  role?: string;
  /** @example "2026-03-25T12:00:41.267Z" */
  createdAt?: string;
}

export interface InternalFeaturesPostTransportHttpCreatePostDTO {
  /**
   * @minLength 1
   * @maxLength 280
   * @example "Hello, world!"
   */
  content: string;
  /** @example "https://example.com/image.jpg" */
  image_url?: string;
  /** @example "123e4567-e89b-12d3-a456-426614174000" */
  parent_id?: string;
}

export interface InternalFeaturesPostTransportHttpLikeResponse {
  message?: string;
}

export interface InternalFeaturesPostTransportHttpPostDTOResponse {
  /** @example "Hello, world!" */
  content?: string;
  /** @example "2026-03-25T12:00:41.267Z" */
  created_at?: string;
  /** @example "3fa85f64-5717-4562-b3fc-2c963f66afa6" */
  id?: string;
  /** @example "https://example.com/image.jpg" */
  image_url?: string;
  /** @example 10 */
  likes_count?: number;
  /** @example "123e4567-e89b-12d3-a456-426614174000" */
  parent_id?: string;
  /** @example "3fa85f64-5717-4562-b3fc-2c963f66afa6" */
  user_id?: string;
}

export interface InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse {
  /** @example "https://example.com/avatar.jpg" */
  avatar_url?: string;
  /** @example "Hello, world!" */
  content?: string;
  /** @example "2026-03-25T12:00:41.267Z" */
  created_at?: string;
  /** @example "2026-03-25T12:00:41.267Z" */
  deleted_at?: string;
  /** @example "test" */
  display_name?: string;
  /** @example "3fa85f64-5717-4562-b3fc-2c963f66afa6" */
  id?: string;
  /** @example "https://example.com/image.jpg" */
  image_url?: string;
  /** @example false */
  is_liked?: boolean;
  /** @example 10 */
  likes_count?: number;
  /** @example "123e4567-e89b-12d3-a456-426614174000" */
  parent_id?: string;
  /** @example "3fa85f64-5717-4562-b3fc-2c963f66afa6" */
  user_id?: string;
  /** @example "testuser" */
  username?: string;
}

export interface InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse {
  pagination?: InternalFeaturesPostTransportHttpPaginationDTO;
  posts?: InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse[];
}

export interface InternalFeaturesPostTransportHttpPaginationDTO {
  /** @example 1 */
  page?: number;
  /** @example 20 */
  page_size?: number;
  /** @example 20 */
  total?: number;
  /** @example 100 */
  total_pages?: number;
}

export interface InternalFeaturesProfileTransportHttpFollowResponse {
  /** @example "success" */
  message?: string;
}

export interface InternalFeaturesProfileTransportHttpProfileResponse {
  /** @example "http://localhost:8080/avatar_url" */
  avatar_url?: string;
  /** @example "2006-01-02T15:04:05Z07:00" */
  created_at?: string;
  /** @example "mark" */
  display_name?: string;
  /** @example "test.user@email.com" */
  email?: string;
  /** @example false */
  email_verified?: boolean;
  /** @example 10 */
  followers_count?: number;
  /** @example 5 */
  following_count?: number;
  /** @example "8794e589-0ddb-43ce-9f92-16faafcf4ee4" */
  id?: string;
  /** @example "user" */
  role?: string;
  /** @example "something smart" */
  status?: string;
  /** @example "username" */
  username?: string;
}

export interface InternalFeaturesProfileTransportHttpUpdateProfileRequest {
  /** @example "mark" */
  display_name?: string;
  /** @example "something smart" */
  status?: string;
  /** @example "username" */
  username?: string;
}
