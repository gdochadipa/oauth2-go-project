CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE IF NOT EXISTS "GrantTypes" AS ENUM (
  'client_credentials',
  'authorization_code',
  'refresh_token',
  'implicit',
  'password'
);

CREATE TYPE IF NOT EXISTS "CodeChallengeMethod" AS ENUM (
  'S256',
  'plain'
);

-- Create User table
CREATE TABLE IF NOT EXISTS  IF NOT EXISTS "User" (
  "id" PRIMARY KEY DEFAULT uuid_generate_v7(),
  "email" VARCHAR UNIQUE NOT NULL,
  "passwordHash" VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS "idx_user_email" ON "User"("email");

-- Create OAuthScope table
CREATE TABLE IF NOT EXISTS  "OAuthScope" (
  "id" PRIMARY KEY DEFAULT uuid_generate_v7(),
  "name" VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS "idx_oauthscope_name" ON "OAuthScope"("name");

-- Create OAuthClient table
CREATE TABLE IF NOT EXISTS  "OAuthClient" (
  "id" PRIMARY KEY DEFAULT uuid_generate_v7(),
  "name" VARCHAR(255) NOT NULL,
  "secret" VARCHAR(255),
  "redirectUris" VARCHAR[] NOT NULL,
  "allowedGrants" "GrantTypes"[] NOT NULL
);

-- Create junction table for OAuthClient and OAuthScope
CREATE TABLE IF NOT EXISTS  "oauthClient_oauthScope" (
  "clientId" UUID NOT NULL,
  "scopeId" UUID NOT NULL,
  PRIMARY KEY ("clientId", "scopeId"),
  FOREIGN KEY ("clientId") REFERENCES "OAuthClient"("id"),
  FOREIGN KEY ("scopeId") REFERENCES "OAuthScope"("id")
);

CREATE INDEX IF NOT EXISTS "idx_oauthclient_oauthscope_clientid" ON "oauthClient_oauthScope"("clientId");
CREATE INDEX IF NOT EXISTS "idx_oauthclient_oauthscope_scopeid" ON "oauthClient_oauthScope"("scopeId");

-- Create OAuthAuthCode table
CREATE TABLE IF NOT EXISTS  "OAuthAuthCode" (
  "code" VARCHAR PRIMARY KEY,
  "redirectUri" VARCHAR,
  "codeChallenge" VARCHAR,
  "codeChallengeMethod" "CodeChallengeMethod" NOT NULL DEFAULT 'plain',
  "expiresAt" TIMESTAMP NOT NULL,
  "userId" UUID,
  "clientId" UUID NOT NULL,
  FOREIGN KEY ("userId") REFERENCES "User"("id"),
  FOREIGN KEY ("clientId") REFERENCES "OAuthClient"("id")
);

-- Create junction table for OAuthAuthCode and OAuthScope
CREATE TABLE IF NOT EXISTS  "OAuthAuthCode_OAuthScope" (
  "oauthAuthCodeCode" VARCHAR NOT NULL,
  "oauthScopeId" UUID NOT NULL,
  PRIMARY KEY ("oauthAuthCodeCode", "oauthScopeId"),
  FOREIGN KEY ("oauthAuthCodeCode") REFERENCES "OAuthAuthCode"("code"),
  FOREIGN KEY ("oauthScopeId") REFERENCES "OAuthScope"("id")
);

-- Create OAuthToken table
CREATE TABLE IF NOT EXISTS  "OAuthToken" (
  "accessToken" VARCHAR PRIMARY KEY,
  "accessTokenExpiresAt" TIMESTAMP NOT NULL,
  "refreshToken" VARCHAR UNIQUE,
  "refreshTokenExpiresAt" TIMESTAMP,
  "clientId" UUID NOT NULL,
  "userId" UUID,
  FOREIGN KEY ("clientId") REFERENCES "OAuthClient"("id"),
  FOREIGN KEY ("userId") REFERENCES "User"("id")
);

CREATE INDEX IF NOT EXISTS "idx_oauthtoken_accesstoken" ON "OAuthToken"("accessToken");
CREATE INDEX IF NOT EXISTS "idx_oauthtoken_refreshtoken" ON "OAuthToken"("refreshToken");

-- Create junction table for OAuthToken and OAuthScope
CREATE TABLE IF NOT EXISTS  "OAuthToken_OAuthScope" (
  "oauthTokenAccessToken" VARCHAR NOT NULL,
  "oauthScopeId" UUID NOT NULL,
  PRIMARY KEY ("oauthTokenAccessToken", "oauthScopeId"),
  FOREIGN KEY ("oauthTokenAccessToken") REFERENCES "OAuthToken"("accessToken"),
  FOREIGN KEY ("oauthScopeId") REFERENCES "OAuthScope"("id")
);