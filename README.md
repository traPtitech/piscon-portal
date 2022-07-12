# PISCON-PORTAL

piscon 用のポータルサイト

## デプロイ

- 対象の過去問の競技鯖 AMI をとってきてベンチマーク鯖として建てる (AMI にベンチマーカーも一緒に入っているため)
- Elastic IP、セキュリティグループなどを設定]
- ssh で入って piscon-portal を `git clone`
- nginx, mariadb 他常駐しているもの (ISUCON11 なら JIA-xxx があった) を停止する
- `make deploy`

で立つはず

## 既知のバグ

セッション周りがバグっていて他のユーザーでログインが発生する場合があります

## 環境変数

直下に`.env`を置いてそこに配置

| key                | value                            |
| ------------------ | -------------------------------- |
| `ENV`              | `prod`or `""`                    |
| `OAUTH_CLIENT_ID`  | traP の OAuth2 クライアントの ID |
| `BENCH_PRIVATE_IP` | ベンチマーカーの Private IP      |
| `MARIADB_USERNAME` | DB のユーザーネーム              |
| `MARIADB_PASSWORD` | DB のパスワード                  |
| `MARIADB_PORT`     | DB のポート番号                  |
| `MARIADB_HOSTNAME` | DB のホスト名                    |
| `MARIADB_DATABASE` | DB のデータベース名              |

### AWS

| key                     | value                   |
| ----------------------- | ----------------------- |
| `AWS_SUBNET_ID`         | VPC のサブネット ID     |
| `AWS_SECURITY_GROUP_ID` | セキュリティグループ ID |
| `AWS_ACCESS_KEY`        | AWS のアクセス ID       |
| `AWS_ACCESS_SECRET`     | AWS のシークレットキー  |

### conoha

メンテナンスされていません

### セットアップ

Docker が入っていることが必要です

## Client

https://github.com/epicmaxco/vuestic-admin

Vuestic admin ベースに構築されています
