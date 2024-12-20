CREATE TABLE "OAuthCode" (
  "id" string,
  "code" string,
  "redirectUri" string,
  "codeChallenge" string,
  "codeChallengeMethod" string,
  "userId" string,
  "expiresAt" datetime,
  "client_id" string,
  "scopes" string[],
  "createdAt" datetime
);

CREATE TABLE "OAuthClient" (
  "id" string,
  "name" string,
  "secret" string,
  "redirectUris" string[],
  "allowedGrants" string[],
  "scopes" string[]
);

CREATE TABLE "OAuthScope" (
  "id" string,
  "name" string
);

CREATE TABLE "OAuthUser" (
  "id" string,
  "username" string,
  "email" string,
  "password" string
);

CREATE TABLE "OAuthToken" (
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
