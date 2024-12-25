CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE code_challenge_method AS ENUM ('S256','plain');

CREATE TYPE grant_allowed_type as ENUM ( 'client_credentials', 'authorization_code', 'refresh_token', 'implicit', 'password')

CREATE TABLE IF NOT EXISTS "OAuthCode" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  "code" VARCHAR(255) UNIQUE NOT NULL,
  "redirectUri" VARCHAR(255),
  "codeChallenge" VARCHAR(255),
  "codeChallengeMethod" code_challenge_method NOT NULL DEFAULT 'plain',
  "userId" UUID NOT NULL REFERENCES OAuthUser(id),
  "expiresAt" datetime,
  "client_id"  UUID NOT NULL REFERENCES OAuthClient(id),
  "scopes" string[],
  "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "OAuthClient" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "secret" VARCHAR(255),
  "redirectUri" VARCHAR(255),
  "allowedGrants" grant_allowed_type[],
  "scopes" string[]
);

CREATE TABLE IF NOT EXISTS "OAuthScope" (
  "id" string,
  "name" string
);

CREATE TABLE IF NOT EXISTS "OAuthUser" (
  "id" string,
  "username" string,
  "email" string,
  "password" string
);

CREATE TABLE IF NOT EXISTS "OAuthToken" (
  "accessToken" string,
  "accessTokenExpiresAt" Datetime,
  "refreshToken" string,
  "refreshTokenExpiresAt" Datetime,
  "clientId" string,
  "userId" string,
  "scopes" string[]
);

ALTER TABLE "OAuthClient" ADD FOREIGN KEY ("id") REFERENCES "OAuthToken" ("clientId");

ALTER TABLE "OAuthUser" ADD FOREIGN KEY ("id") REFERENCES "OAuthCode" ("userId");

ALTER TABLE "OAuthUser" ADD FOREIGN KEY ("id") REFERENCES "OAuthToken" ("userId");

ALTER TABLE "OAuthScope" ADD FOREIGN KEY ("id") REFERENCES "OAuthCode" ("scopes");

ALTER TABLE "OAuthScope" ADD FOREIGN KEY ("id") REFERENCES "OAuthClient" ("scopes");

ALTER TABLE "OAuthScope" ADD FOREIGN KEY ("id") REFERENCES "OAuthToken" ("scopes");
