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
  InternalFeaturesAuthTransportHttpRegisterDTO,
  InternalFeaturesAuthTransportHttpRegisterDTOResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Register<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Создаёт нового пользователя и возвращает его данные
   *
   * @tags Auth
   * @name RegisterCreate
   * @summary Регистрация пользователя
   * @request POST:/register
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
      path: `auth/register`,
      method: "POST",
      body: request,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
