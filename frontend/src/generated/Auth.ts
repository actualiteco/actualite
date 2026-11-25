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
  InternalFeaturesAuthTransportHttpLoginDTO,
  InternalFeaturesAuthTransportHttpRefreshDTO,
  InternalFeaturesAuthTransportHttpRegisterDTO,
  InternalFeaturesAuthTransportHttpRegisterDTOResponse,
  InternalFeaturesAuthTransportHttpTokensDTOResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Auth<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Авторизует пользователя по email и паролю, возвращает access и refresh токены с временем жизни access token в секундах (10 минут)
   *
   * @tags Auth
   * @name LoginCreate
   * @summary Авторизация по email и паролю
   * @request POST:/auth/login
   */
  loginCreate = (
    request: InternalFeaturesAuthTransportHttpLoginDTO,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesAuthTransportHttpTokensDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/auth/login`,
      method: "POST",
      body: request,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Инвалидирует refresh-токен. Пока привязанный к нему access-токен не протух, по нему еще можно будет делать запросы (будет исправлено позже после интеграции Redis), но при это обновлять refresh-токен уже нельзя. Операция идемпотентна, можно несколько раз отозвать один и тот же refresh-токен без ошибки
   *
   * @tags Auth
   * @name LogoutCreate
   * @summary Инвалидация refresh-токена
   * @request POST:/auth/logout
   * @secure
   */
  logoutCreate = (params: RequestParams = {}) =>
    this.request<
      void,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/auth/logout`,
      method: "POST",
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Обновляет refresh-токен, возвращает access и refresh токены с временем жизни access token в секундах (10 минут). Используется, когда прежний AccessToken протух
   *
   * @tags Auth
   * @name RefreshCreate
   * @summary Обновление refresh-токена
   * @request POST:/auth/refresh
   */
  refreshCreate = (
    request: InternalFeaturesAuthTransportHttpRefreshDTO,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesAuthTransportHttpTokensDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/auth/refresh`,
      method: "POST",
      body: request,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Создаёт нового пользователя и возвращает его данные
   *
   * @tags Auth
   * @name RegisterCreate
   * @summary Регистрация пользователя
   * @request POST:/auth/register
   */
  registerCreate = (
    request: InternalFeaturesAuthTransportHttpRegisterDTO,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesAuthTransportHttpRegisterDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `/auth/register`,
      method: "POST",
      body: request,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
