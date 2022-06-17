# PISCON-PORTAL

piscon用のポータルサイト


## 既知のバグ

セッション周りがバグっていて他のユーザーでログインが発生する場合があります

## 環境変数
直下に`.env`を置いてそこに配置


| key                 | value                |
|---------------------|----------------------|
| `ENV`               | `prod`or `""`        |
| `OAUTH_CLIENT_ID`   | traPのOAuth2クライアントのID |
| `ACCESS_ID`         | AWSのアクセスID           |
| `ACCESS_SECRET_KEY` | AWSのシークレットキー         |
| `MARIADB_USERNAME`  | DBのユーザーネーム           |
| `MARIADB_PASSWORD`  | DBのパスワード             |
| `MARIADB_PORT`      | DBのポート番号             |
| `MARIADB_HOSTNAME`  | DBのホスト名              |
| `MARIADB_DATABASE`  | DBのデータベース名           |

### AWS
| key                     | value        |
|-------------------------|--------------|
| `AWS_SUBNET_ID`         | VPCのサブネットID  |
| `AWS_SECURITY_GROUP_ID` | セキュリティグループID |
| `AWS_ACCESS_KEY`        | AWSのアクセスID   |
| `AWS_ACCESS_SECRET`     | AWSのシークレットキー |


### conoha
メンテナンスされていません


### セットアップ

Dockerが入っていることが必要です


## Client

https://github.com/epicmaxco/vuestic-admin

Vuestic admin ベースに構築されています
