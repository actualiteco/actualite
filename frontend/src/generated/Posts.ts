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

import {
  ActualitecoActualiteInternalCoreDomainCustomError,
  ActualitecoActualiteInternalCoreDomainInternalError,
  InternalFeaturesPostTransportHttpCreatePostDTO,
  InternalFeaturesPostTransportHttpImageUrlDTOResponse,
  InternalFeaturesPostTransportHttpLikeResponse,
  InternalFeaturesPostTransportHttpPostDTOResponse,
  InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse,
  InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Posts<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Возвращаются посты подписок пользователя и посты за последний месяц (отсортированные по количеству лайков) поддержкой пагинации через параметры page и page_size
   *
   * @tags Posts
   * @name PostsList
   * @summary Возвращает посты для ленты
   * @request GET:/posts
   * @secure
   */
  postsList = (
    query?: {
      /**
       * Номер страницы
       * @min 1
       * @default 1
       */
      page?: number;
      /**
       * Размер страницы
       * @min 1
       * @max 30
       * @default 20
       */
      page_size?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Создает пост длиной <= 280 символов
   *
   * @tags Posts
   * @name CreateCreate
   * @summary Создание поста
   * @request POST:/posts/create
   * @secure
   */
  createCreate = (
    request: InternalFeaturesPostTransportHttpCreatePostDTO,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesPostTransportHttpPostDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/create`,
      method: "POST",
      body: request,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Загрузить изображение для поста (jpg или png, до 5 МБ) в хранилище
   *
   * @tags Posts
   * @name ImageCreate
   * @summary Загрузить изображение
   * @request POST:/posts/image
   * @secure
   */
  imageCreate = (
    data: {
      /**
       * Файл изображения (jpg/png, до 5 МБ)
       * @format binary
       */
      file: File;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesPostTransportHttpImageUrlDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/image`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.FormData,
      format: "json",
      ...params,
    });
  /**
   * @description Ищет посты пользователя по его Access Token с поддержкой пагинации через параметры page и page_size
   *
   * @tags Posts
   * @name GetPosts
   * @summary Поиск постов пользователя по его JWT
   * @request GET:/posts/me
   * @secure
   */
  getPosts = (
    query?: {
      /**
       * Номер страницы
       * @min 1
       * @default 1
       */
      page?: number;
      /**
       * Размер страницы
       * @min 1
       * @max 30
       * @default 20
       */
      page_size?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/me`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Ищет пост по ID и возвращает информацию о посте + краткую информацию об авторе (username, display_name, avatar_url)
   *
   * @tags Posts
   * @name PostsDetail
   * @summary Поиск поста по ID
   * @request GET:/posts/{postID}
   * @secure
   */
  postsDetail = (postId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesPostTransportHttpPostWithAuthorDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/${postId}`,
      method: "GET",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Удаляет пост, если запрос на удаление исходит от автора поста. Метод идемпотентен
   *
   * @tags Posts
   * @name PostsDelete
   * @summary Удаление поста
   * @request DELETE:/posts/{postID}
   * @secure
   */
  postsDelete = (postId: string, params: RequestParams = {}) =>
    this.request<
      void,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/${postId}`,
      method: "DELETE",
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Поставить лайк на пост авторизованному пользователю
   *
   * @tags Posts
   * @name LikeCreate
   * @summary Поставить лайк
   * @request POST:/posts/{postID}/like
   * @secure
   */
  likeCreate = (postId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesPostTransportHttpLikeResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/${postId}/like`,
      method: "POST",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Ищет комментарии на пост по его ID с поддержкой пагинации через параметры page и page_size
   *
   * @tags Posts
   * @name RepliesList
   * @summary Поиск комментариев на пост
   * @request GET:/posts/{postID}/replies
   * @secure
   */
  repliesList = (
    postId: string,
    query?: {
      /**
       * Номер страницы
       * @min 1
       * @default 1
       */
      page?: number;
      /**
       * Размер страницы
       * @min 1
       * @max 30
       * @default 20
       */
      page_size?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/${postId}/replies`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Убрать лайк с поста авторизованному пользователю
   *
   * @tags Posts
   * @name UnlikeCreate
   * @summary Убрать лайк
   * @request POST:/posts/{postID}/unlike
   * @secure
   */
  unlikeCreate = (postId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesPostTransportHttpLikeResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/posts/${postId}/unlike`,
      method: "POST",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
