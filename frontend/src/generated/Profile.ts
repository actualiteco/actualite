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
  InternalFeaturesPostTransportHttpPostsWithPaginationDTOResponse,
  InternalFeaturesProfileTransportHttpFollowResponse,
  InternalFeaturesProfileTransportHttpProfileResponse,
  InternalFeaturesProfileTransportHttpUpdateProfileRequest,
  InternalFeaturesProfileTransportHttpUploadAvatarResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Profile<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Обновить профиль пользователя (display_name, username, status)
   *
   * @tags Profile
   * @name ProfilePartialUpdate
   * @summary Обновить профиль пользователя
   * @request PATCH:/profile
   * @secure
   */
  profilePartialUpdate = (
    request: InternalFeaturesProfileTransportHttpUpdateProfileRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesProfileTransportHttpProfileResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile`,
      method: "PATCH",
      body: request,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Загрузить аватар пользователя (jpg или png, до 2 МБ) в хранилище и обновить профиль
   *
   * @tags Profile
   * @name AvatarCreate
   * @summary Загрузить аватар
   * @request POST:/profile/avatar
   * @secure
   */
  avatarCreate = (
    data: {
      /**
       * Файл аватара (jpg/png, до 2 МБ)
       * @format binary
       */
      file: File;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesProfileTransportHttpUploadAvatarResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile/avatar`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.FormData,
      format: "json",
      ...params,
    });
  /**
   * @description Получить профиль пользователя с постами по JWT
   *
   * @tags Profile
   * @name GetProfile
   * @summary Получить профиль пользователя по JWT
   * @request GET:/profile/me
   * @secure
   */
  getProfile = (params: RequestParams = {}) =>
    this.request<
      InternalFeaturesProfileTransportHttpProfileResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile/me`,
      method: "GET",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Получить профиль пользователя с постами по id пользователя
   *
   * @tags Profile
   * @name ProfileDetail
   * @summary Получить профиль пользователя
   * @request GET:/profile/{userID}
   * @secure
   */
  profileDetail = (userId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesProfileTransportHttpProfileResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile/${userId}`,
      method: "GET",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Подписаться авторизованным пользователем на другого пользователя по userID
   *
   * @tags Profile
   * @name FollowCreate
   * @summary Подписаться на пользователя
   * @request POST:/profile/{userID}/follow
   * @secure
   */
  followCreate = (userId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesProfileTransportHttpFollowResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile/${userId}/follow`,
      method: "POST",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Ищет посты пользователя по ID с поддержкой пагинации через параметры page и page_size
   *
   * @tags Posts
   * @name PostsList
   * @summary Поиск постов пользователя
   * @request GET:/profile/{userID}/posts
   * @secure
   */
  postsList = (
    userId: string,
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
      path: `/profile/${userId}/posts`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Отписаться авторизованным пользователем от другого пользователя по userID
   *
   * @tags Profile
   * @name UnfollowCreate
   * @summary Отписаться от пользователя
   * @request POST:/profile/{userID}/unfollow
   * @secure
   */
  unfollowCreate = (userId: string, params: RequestParams = {}) =>
    this.request<
      InternalFeaturesProfileTransportHttpFollowResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/profile/${userId}/unfollow`,
      method: "POST",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
