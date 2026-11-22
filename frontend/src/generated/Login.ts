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
  InternalFeaturesAuthTransportHttpLoginDTOResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Login<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Авторизует пользователя по email и паролю, возвращает JWT
   *
   * @tags Auth
   * @name LoginCreate
   * @summary Авторизация по email и паролю
   * @request POST:/login
   */
  loginCreate = (
    request: InternalFeaturesAuthTransportHttpLoginDTO,
    params: RequestParams = {},
  ) =>
    this.request<
      InternalFeaturesAuthTransportHttpLoginDTOResponse,
      | ActualitecoActualiteInternalCoreDomainCustomError
      | ActualitecoActualiteInternalCoreDomainInternalError
    >({
      path: `auth/login`,
      method: "POST",
      body: request,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
